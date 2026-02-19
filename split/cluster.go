package split

import (
	"fmt"
	"math"
	"strings"

	"github.com/ademajagon/gix/internal/git"
	"github.com/ademajagon/gix/provider"
)

type HunkGroup struct {
	Hunks   []git.Hunk
	Message string
}

const similarityThreshold = 0.85

// ClusterHunks uses embedding-based cosine similarity to group hunks than generates commit message for each group
func ClusterHunks(p provider.AIProvider, hunks []git.Hunk) ([]HunkGroup, error) {
	if len(hunks) == 0 {
		return nil, nil
	}

	texts := make([]string, len(hunks))
	for i, h := range hunks {
		texts[i] = h.FilePath + "\n" + h.Header + "\n" + h.Body
	}

	embeddings, err := p.GetEmbeddings(texts)
	if err != nil {
		return nil, fmt.Errorf("embedding hunks: %w", err)
	}

	if len(embeddings) != len(hunks) {
		return nil, fmt.Errorf("embedding count mismatch: got %d, want %d", len(embeddings), len(hunks))
	}

	used := make([]bool, len(hunks))
	var groups []HunkGroup

	for i := range hunks {
		if used[i] {
			continue
		}
		group := []git.Hunk{hunks[i]}
		used[i] = true

		for j := i + 1; j < len(hunks); j++ {
			if used[j] {
				continue
			}
			if cosineSimilarity(embeddings[i], embeddings[j]) >= similarityThreshold {
				group = append(group, hunks[j])
				used[j] = true
			}
		}

		groups = append(groups, HunkGroup{Hunks: group})
	}

	for i := range groups {
		patch := joinPatch(groups[i].Hunks)
		msg, err := p.GenerateCommitMessage(patch)
		if err != nil {
			return nil, fmt.Errorf("generating message for group %d: %w", i+1, err)
		}
		groups[i].Message = msg
	}

	return groups, nil
}

func joinPatch(hunks []git.Hunk) string {
	var b strings.Builder
	for _, h := range hunks {
		b.WriteString(h.Body)
		b.WriteString("\n\n")
	}
	return b.String()
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		ai, bi := float64(a[i]), float64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
