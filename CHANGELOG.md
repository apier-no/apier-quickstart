# Changelog

All notable changes to apier-quickstart are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `examples/curl`, `examples/typescript`, `examples/python` and `examples/go`:
  Verify, Authority and Proof examples with no dependencies. Verify and
  Authority run without a key against the public sandbox (org 999999999).
- `compare/README.md`: what Brreg's open APIs return for signing authority
  and what the Apier Authority response returns.
- `.github/workflows/sandbox-smoke.yml`: runs the examples against the live
  sandbox on push, on pull request and weekly.

### Changed
- README rebuilt around company verification and signing authority. The
  Altinn 2 to Altinn 3 lookup moved to "Other endpoints".

### Removed
- README references to `apier_test_` keys, a "free tier" call allowance,
  the `X-Apier-Key` header and a Maskinporten guide. None of them describes
  the current service.

## [0.2.0] — 2026-05-22

### Added
- "Get a key in 30 seconds" onboarding section in README
- APIER_BASE_URL self-hosted / staging documentation
- Working `npx @apier-no/mcp` one-shot MCP server example
- Dependabot config (.github/dependabot.yml) — weekly npm + pip +
  github-actions
- `.venv/`, `dist/`, `build/` added to `.gitignore`
- CHANGELOG.md (this file)

### Changed
- Stale `@apier/mcp` references corrected to `@apier-no/mcp` (the
  published npm package name as of PR-083 v0.1.1, 2026-05-20)
- June 19, 2026 Altinn 2 shutdown deadline elevated to above-the-fold
  urgency framing

### Fixed
- Removed "Once PR-083 ships" TODO — PR-083 shipped 2026-05-21,
  working `npx @apier-no/mcp` example now in place

### Security
- Examples reject non-HTTPS `APIER_BASE_URL` (allowing only `https://`
  or `http://localhost`) to prevent accidental key exfiltration when
  the env var is set to an untrusted host.

## [0.1.0] — 2026-05-11

### Added
- Initial scaffold: Node.js + Python examples hitting the Altinn
  migration bridge at /api/v1/tools/altinn-migration
- GitHub Actions sanity-check workflow with positive-control assertion
- MIT LICENSE
- SECURITY.md → security@apier.no
