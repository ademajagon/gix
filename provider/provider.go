package provider

const (
	// CommitMessageSystemPrompt is the system instruction for generating commit messages
	CommitMessageSystemPrompt = "You are a concise assistant that only returns a one-line, conventional commit message. No explanations, markdown, or commentary."

	// CommitMessageUserPromptTemplate is the template for the user prompt when generating commit messages
	CommitMessageUserPromptTemplate = "Write a single-line conventional commit message that describes the following Git diff. Only return the commit message. Do not include explanations, newlines, or formatting beyond the message itself. Diff:\n\n"
)

// AIProvider abstracts chat completion and embedding capabilities
// so that different backends (OpenAI, Gemini, etc.) can be used interchangeably.
type AIProvider interface {
	GenerateCommitMessage(diff string) (string, error)
	GetEmbeddings(texts []string) ([][]float32, error)
}
