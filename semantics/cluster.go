package semantics

import (
	"fmt"
	"math"
	"strings"

	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/openai"
)

type HunkGroup struct {
	Hunks   []git.Hunk
	Message string
}

func ClusterHunks(apiKey string, hunks []git.Hunk) ([]HunkGroup, error) {
	if len(hunks) == 0 {
		return nil, nil
	}

	var texts []string
	for _, h := range hunks {
		texts = append(texts, fmt.Sprintf("%s\n%s", h.Header, h.Body))
	}

	embeddings, err := openai.GetEmbeddings(apiKey, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to embed hunks: %w", err)
	}

	const threshold = 0.85
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

			if cosineSimilarity(embeddings[i], embeddings[j]) >= threshold {
				group = append(group, hunks[j])
				used[j] = true
			}

		}

		groups = append(groups, HunkGroup{
			Hunks: group,
		})
	}

	for i := range groups {
		diff := joinGroupPatch(groups[i].Hunks)
		msg, err := openai.GenerateCommitMessage(apiKey, diff)
		if err != nil {
			return nil, fmt.Errorf("failed to generate message for group %d: %w", i+1, err)
		}
		groups[i].Message = msg
	}

	return groups, nil
}

func joinGroupPatch(hunks []git.Hunk) string {
	var b strings.Builder
	for _, h := range hunks {
		b.WriteString("diff --git a/" + h.FilePath + " b/" + h.FilePath + "\n")
		b.WriteString(h.Body + "\n\n")
	}

	return b.String()
}

func cosineSimilarity(a, b []float32) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
