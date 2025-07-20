<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/logo_gix_light.svg">
  <img alt="gix logo" src="/docs/logo_gix_dark.svg" width="50%" height="50%">
</picture>

gix: Git on the command line, with a bit of AI.

[![Release](https://img.shields.io/github/v/release/ademajagon/gix?color=green&label=release)](https://github.com/ademajagon/gix/releases)

</div>

---

## Overview

Gix is a CLI tool that helps you keep your git history clean. It can split large diffs, write conventional commits, and automate the repetitive parts.

It runs locally, uses your own OpenAI key, and fits into your existing workflow.

---

## Features

- AI-suggested conventional commit messages
- `gix split` - split staged diffs into multiple commits
- Groups related changes using LLM-based embeddings
- Bring your own OpenAI API key (no lock-in)
- Built in Go – fast, portable, and cross-platform

---

## Installation

### macOS (Homebrew)

```bash
brew tap ademajagon/gix
brew install gix
```

### Linux / Windows

Download from [Releases](https://github.com/ademajagon/gix/releases) and add it to your `PATH`.


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

### Split staged changes into atomic commits

```bash
git add -p
gix split
```

Gix will group commits and ask for confirmation before applying.

---

## Configuration

```bash
gix config set openai_key sk-...
```


## License

MIT © [Agon Ademaj](https://github.com/ademajagon)