## [v0.3.0] - 2026-03-01

### Added
- Ollama provider support, run `gix commit` and `gix split` fully offline with local models
- `gix config set-provider ollama` no API key required
- `gix config set-ollama-url <url>` point gix at a remote Ollama instance
- `gix config set-ollama-model <chat-model> <embed-model>` to set specific models

### Changed
- Shared HTTP client extracted into `chatclient.go`, OpenAI, Ollama and other LLMs share the same transport layer with OpenAI-compatible API
- Commit message prompt refactor to improve output quality across smaller local models
- README updated with Ollama setup instructions and cleaner structure

### Fixed
- `gix config set-provider` now accepts `ollama` as a valid value
- Early API key guard in registry blocked Ollama before reaching the provider switch

## [v0.2.9] - 2026-02-23
### Added

- Automated Homebrew tap formula update on every release, `brew upgrade gix` now works automatically
- checksums.txt attached to every GitHub Release for binary verification

### Changed
- Release pipeline split into build and release jobs so checksums are generated after all binaries are ready

## [v0.2.8] - 2026-02-22
### Fixed
Update notice now shows even when a command fails

## [v0.2.8] - 2026-02-22
### Fixed
- Version update notice now shows even when a command fails

## [v0.2.7] - 2026-02-22

### Added
- Checkpoint system for version update notifications (mirrors Terraform's go-checkpoint pattern)
- `gix config update-check off` to disable update checks permanently
- `GIX_CHECKPOINT_DISABLE=1` env var for one-off session disable
- Gemini provider support
- `gix split` beta command

### Fixed
- Hunk.Header was always empty in gix split
- Config save no longer prints debug output