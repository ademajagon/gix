## [v0.2.0] - 2026-02-22

### Added
- Checkpoint system for version update notifications (mirrors Terraform's go-checkpoint pattern)
- `gix config update-check off` to disable update checks permanently
- `GIX_CHECKPOINT_DISABLE=1` env var for one-off session disable
- Gemini provider support
- `gix split` beta command

### Fixed
- Hunk.Header was always empty in gix split
- Config save no longer prints debug output