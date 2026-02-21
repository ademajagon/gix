<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/logo_gix_light.svg">
  <img alt="gix logo" src="/docs/logo_gix_dark.svg" width="50%" height="50%">
</picture>

gix: Git on the command line, with a bit of AI.

[![Release](https://img.shields.io/github/v/release/ademajagon/gix?color=green&label=release)](https://github.com/ademajagon/gix/releases)
[![Build](https://github.com/ademajagon/gix/actions/workflows/ci.yml/badge.svg)](https://github.com/ademajagon/gix/actions/workflows/ci.yml)

</div>

---

## Overview

Gix is a CLI tool that helps you keep your git history clean. It can write conventional commits, split large diffs, and automate the repetitive git parts.

It runs locally, uses your own API key (OpenAI or Gemini), and fits into your existing workflow.

---

## Features

- AI-suggested conventional commit messages
- `gix split` - split staged diffs into multiple commits
- Groups related changes using LLM-based embeddings
- **Multiple AI providers** - OpenAI or Google Gemini
- Bring your own API key (no lock-in)
- Built in Go - fast, portable, and cross-platform

---

## Installation

### macOS (Homebrew)

```bash
brew tap ademajagon/gix
brew install gix
```

### Linux / Windows

Download binaries from [Releases](https://github.com/ademajagon/gix/releases) and add it to your `PATH`.

### Go (for contributors)

```bash
go install github.com/ademajagon/gix@latest
```

## Usage

### Generate a commit message

```bash
git add .
gix commit
```

### Split staged changes (beta)

```bash
git add .
gix split
```

Gix will group commits and ask for confirmation before applying.

---

## Configuration

`gix` stores configuration locally on your machine.

### Set provider
```bash
gix config set-provider openai
gix config set-provider gemini
```

### Set API key

```bash
# open-ai (default)
gix config set-key

# gemini
gix config set-key --provider gemini
```

Configure both providers and switch anytime.

---

## Supported Providers

| Provider | Chat Model          | Embeddings             |
| -------- | ------------------- |------------------------|
| OpenAI   | gpt-4o              | text-embedding-3-small |
| Gemini   | gemini-flash-latest | gemini-embedding-001   |

---

## License

MIT © [Agon Ademaj](https://github.com/ademajagon)
