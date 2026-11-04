# Owner-curated audit stage queue

- **Only staged:** `gen-type-go-pr6` (`STAGED` + `enabled=true`). It retains the SHA-locked original PR #6 replay from source branch commit `f046b362a482d1fa4d557df284db3cf39bf0a55d`. Running the master setup never calls its replay.
- **Historical, inactive:** `kol-merkle-pr4` (`REFERENCE_ONLY` + `enabled=false`). Its existing `main` scripts and prompt are intentionally unchanged.
- **Queue policy:** explicit selection, exactly one active stage, no auto-runs, no memory capture, no inferred permissions.

Any new profile should be proposed as a **reviewable Git change** that includes exact target SHA/tree, a bounded profile, immutable source/provenance hashes, profile-specific scope and permitted tests, explicit authority flags all set false, and a deterministic runner that uses isolated target directories. A Go validator update must be separately reviewed before adding multiple simultaneous active stages.

Go can emit context via `go run ./tools/masteraudit/main.go -profile gen-type-go-pr6 context` without touching a target. The shell wrapper `sh scripts/stage-audit.sh gen-type-go-pr6` writes an ignored, hash-identified local handoff under `audit-artifacts/`.
