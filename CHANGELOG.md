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