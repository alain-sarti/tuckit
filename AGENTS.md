## Checks

Run this before calling a change done. It runs the same checks as CI and stops at the first failure:

```sh
golangci-lint run && go build ./... && go test ./...
```

- `golangci-lint run` lints and checks formatting. It runs gofmt, go vet, errcheck, staticcheck, ineffassign and unused, as configured in `.golangci.yml`. It needs golangci-lint v2. CI pins v2.15.0.
- `go build` typechecks and compiles every package.
- `go test` runs all tests.

If you change this command, change `.github/workflows/ci.yml` to match, and the other way round.

## Agent skills

### Issue tracker

Issues live in GitHub Issues on alain-sarti/tuckit, managed with `gh`. See `docs/agents/issue-tracker.md`.

### Triage labels

Uses the five default labels: needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one `GLOSSARY.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`.
