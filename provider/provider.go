package provider

// AIProvider is the core abstraction for AI providers.
type AIProvider interface {
	// GenerateCommitMessage returns a Conventional Commit message
	GenerateCommitMessage(diff string) (string, error)

	// GetEmbeddings returns vector embeddings for each input string
	GetEmbeddings(texts []string) ([][]float32, error)
}

const (
	CommitMessageSystem = "You are a concise assistant that only returns a one-line, " +
		"conventional commit message (e.g. feat: add login endpoint). " +
		"No explanations, markdown, code fences, or commentary, only the message."

	CommitMessageUser = "Write a single-line conventional commit message for the following " +
		"git diff. Return only the commit message, nothing else.\n\nDiff:\n\n"
)
