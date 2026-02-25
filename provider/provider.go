package provider

// AIProvider is the core abstraction for AI providers.
type AIProvider interface {
	GenerateCommitMessage(diff string) (string, error)
	GetEmbeddings(texts []string) ([][]float32, error)
}

const CommitMessageSystem = "You are a conventional commit message generator. You only output commit messages, nothing else."

const CommitMessageUser = `You must output ONLY a single conventional commit message. No explanations. No descriptions. No extra text.

Format: <type>(<optional scope>): <description>

Types and when to use them:
- feat: a new feature or capability was introduced
- fix: a bug or incorrect behavior was corrected
- refactor: code was restructured or cleaned up without changing behavior
- chore: maintenance, dependency updates, or tooling changes
- docs: documentation was added or updated
- style: formatting or whitespace changes only, no logic change
- perf: a change that improves performance
- test: tests were added, updated, or fixed
- build: changes to the build system, Makefile, or compilation
- ci: changes to CI/CD pipelines or workflows
- revert: a previous commit was undone

Scope is optional but should reflect the area of the codebase changed (e.g. provider, config, git, cmd).

Examples:
feat(provider): add ollama local inference support
fix(provider): skip api key validation for ollama
refactor(git): truncate large diffs before sending to ai
chore: upgrade go version to 1.24
docs: add ollama setup instructions to readme

Diff:
%s

Conventional commit message:`
