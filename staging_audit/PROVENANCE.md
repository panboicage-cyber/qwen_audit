# Frozen branch provenance (not a proof of freshness)

- Baseline main `af6609fae097efe80742025b0b5bf27b9acfccea`.
- RepoMemory integration `feat/auditgo-curated-memory-v0.1.0` @ `57a08f95db8f297136db6e2a5eed8730d0df6fc1`.
- Gen Type-Go PR6 staging `feat/gen-type-go-pr6-deterministic-audit` @ `f046b362a482d1fa4d557df284db3cf39bf0a55d`.
- PR6 lock Git blob `901c8cf5f4636351ae228ef351b61d97a0d4589b`.
- PR6 replay Git blob `7b225092a877f33f1e582969c9988f13c3687a99`.
- PR6 task prompt Git blob `44cb3dd625c0bf932150ef926333103a31a0b951`.

These identities were observed from GitHub during composition. The master workspace retains those bytes; it does not use their historical review/CI outcomes as fresh evidence. Source head lock `0af4707d66bb30aef09fe80316c408568e256e53` is immutable until a separately reviewed owner update. The inactive `kol-merkle-pr4` stage remains part of the historical catalog only.
