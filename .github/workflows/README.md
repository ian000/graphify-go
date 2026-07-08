# Workflow Notes

- `release.yml` is the trusted-publisher workflow bound in npm package settings.
- Do not rename `release.yml` without updating npm Trusted Publisher configuration first.
- `ci.yml` is only for push/PR verification and is safe to evolve independently.
