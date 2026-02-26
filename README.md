<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/logo_gix_light.svg">
  <img alt="gix logo" src="/docs/logo_gix_dark.svg" width="40%" height="40%">
</picture>

gix: Git on the command line, with a bit of AI.

[![Release](https://img.shields.io/github/v/release/ademajagon/gix?color=green&label=release)](https://github.com/ademajagon/gix/releases)
[![Build](https://github.com/ademajagon/gix/actions/workflows/ci.yml/badge.svg)](https://github.com/ademajagon/gix/actions/workflows/ci.yml)

</div>

---

`gix` generates conventional commit messages from your staged diff and splits large changes into small commits.

## Why gix

Writing a good commit message takes discipline and it's usually the first thing skipped when moving fast. `gix` handles it for you, read the diff, generate the message, commit.

For larger changes, `gix split` groups related hunks by semantic similarity and proposes one atomic commit per group, each with its own generated message. It turns one big noisy commit into a clean history.

---

## Features

- Generates [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) from staged diffs
- `gix split` splits a large diff into multiple semantic commits
- Embedding-based hunk clustering for intelligent grouping
- Multiple AI providers, OpenAI, Gemini or Ollama (local)
- Bring your own API key
- Runs fully offline with Ollama
- Single binary, no runtime dependencies
---

## Installation

### macOS (Homebrew)

```bash
brew tap ademajagon/gix
brew install gix
```

### Linux / Windows

Download binaries from [Releases](https://github.com/ademajagon/gix/releases) and add it to your `PATH`.

### From source

```bash
go install github.com/ademajagon/gix@latest
```

## Usage

### Generate a commit message

```bash
git stage .
gix commit
```

You'll see a suggested message and can accept, edit, regenerate or cancel.

### Split a large diff into multiple commits (beta)

```bash
git stage .
gix split
```
gix analyses the staged diff, groups related hunks and proposes one commit per group

---

## Configuration

`gix` stores configuration locally on your machine.

### Set provider
```bash
gix config set-provider openai    # default
gix config set-provider gemini
gix config set-provider ollama    # local, no API key required
```

### Set API key

```bash
gix config set-key                        # OpenAI
gix config set-key --provider gemini      # Gemini
```

Configure both providers and switch anytime.

---

## Supported Providers

| Provider | Chat Model          | Embeddings             |
| -------- | ------------------- |------------------------|
| OpenAI   | gpt-4o              | text-embedding-3-small |
| Gemini   | gemini-flash-latest | gemini-embedding-001   |
| Ollama   | llama3.1:8b (configurable) | nomic-embed-text (configurable)   |

---

### Update checks

gix checks for new releases in the background after each command and prints a notice if one is available. Results are cached for 48 hours and the check never blocks the primary command.

To disable:
```bash
gix config update-check off
# or for a single session:
GIX_CHECKPOINT_DISABLE=1 gix commit
```

---

## License

MIT © [Agon Ademaj](https://github.com/ademajagon)
