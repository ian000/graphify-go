# Releasing

This repository ships both GitHub Release binaries and an npm wrapper package.

## Defaults

- Source of truth: Git tag `vX.Y.Z`
- Binary distribution: GitHub Releases assets
- npm publish path: GitHub Actions only, via npm Trusted Publisher and OIDC
- Local `npm publish`: do not use
- Changelog source: auto-generated GitHub Release notes
- Package version must exactly match the tag version

## Versioning Policy

- `patch`: bug fix, packaging fix, binary download/install fix
- `minor`: backward-compatible CLI capability, new supported output, new platform-safe behavior
- `major`: breaking CLI flags, output format, install contract, or binary compatibility expectations

## Release Steps

1. Make sure `package.json` version is correct.
2. Run `npm run release:check`.
3. Merge to `main`.
4. Create and push the tag:

```bash
git tag -a v1.1.2 -m "Release v1.1.2"
git push origin v1.1.2
```

5. Watch `.github/workflows/release.yml`.
6. Confirm three outcomes:
   - GitHub Release created
   - binary assets uploaded
   - npm version published

## Changelog Convention

- Use Conventional Commit style when possible: `feat:`, `fix:`, `chore:`, `docs:`
- GitHub Release notes are the public changelog
- No manual `CHANGELOG.md` is required by default

## Rollback Policy

- If preflight or build fails before publish, fix the problem and rerun the workflow.
- If GitHub Release exists but npm publish failed, keep the release record and publish a new patch version.
- If npm publish already succeeded, do not overwrite or republish the same version.
- For binary or installer mistakes after release, cut a new patch tag and publish again.

## Trusted Publisher Guardrails

- Keep the workflow filename stable: `release.yml`
- Keep publishing on GitHub-hosted runners only
- Keep `id-token: write` on the npm publish job
- Do not reintroduce `NPM_TOKEN` secrets for standard releases
