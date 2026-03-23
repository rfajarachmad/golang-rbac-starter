---
name: verify
description: Run full code quality verification — build, vet, lint (zero warnings), tests with race detection, coverage enforcement, and migration consistency check
user_invocable: true
---

# Verify Code Quality Skill

Run a comprehensive quality gate on the codebase. Every check must pass — report results as a checklist and stop at the first failure category so the user can fix before re-running.

## Checks to Run (in order)

Execute each step below. After ALL steps complete, present a summary table. If any step fails, still run the remaining steps but mark the overall result as FAILED.

### 1. Build

```bash
go build ./...
```

Must exit 0. Catches syntax errors, missing imports, type mismatches.

### 2. Go Vet

```bash
go vet ./...
```

Must exit 0 with zero output. Catches suspicious constructs: printf format mismatches, unreachable code, struct tag errors, copying locks.

### 3. Lint — Zero Warnings

```bash
golangci-lint run ./...
```

Must exit 0 with **zero warnings**. The project's `.golangci.yml` already configures 18 linters (errcheck, gosimple, govet, staticcheck, unused, bodyclose, goconst, gocritic, gofmt, goimports, gosec, misspell, prealloc, revive, unconvert, unparam, whitespace, ineffassign).

If there are warnings, list each one with file:line and the fix. Apply the fixes, then re-run lint to confirm zero warnings before proceeding.

### 4. Format Check

```bash
gofmt -l ./internal/ ./cmd/ ./test/
```

Must produce no output (no unformatted files). If files are listed, run `go fmt ./...` to fix them, then re-check.

### 5. Module Tidiness

```bash
go mod tidy
git diff --exit-code go.mod go.sum
```

After `go mod tidy`, go.mod and go.sum must have no changes. If they differ, the module wasn't tidy — report which dependencies changed.

### 6. Tests — All Passing

```bash
go test -v -count=1 ./test/ 2>&1 | cat
```

Must exit 0 with all tests PASS. Report total test count and any failures with their names.

### 7. Race Detection

```bash
go test -race -count=1 ./test/ 2>&1 | cat
```

Must exit 0 with no data race detected. Race conditions are non-deterministic bugs that can corrupt data in production.

### 8. Coverage Enforcement

```bash
./scripts/check-coverage.sh
```

Must pass all thresholds:
- **Overall** ≥ 70%
- **UseCase** ≥ 60%
- **Delivery** ≥ 70%

Report actual vs required for each. If a threshold fails, identify the specific functions with 0% or low coverage and suggest which test cases would cover them.

### 9. Migration Consistency

Check that:
- Every `.up.sql` has a matching `.down.sql` in `db/migrations/`
- Sequence numbers are contiguous (no gaps)
- No duplicate sequence numbers

```bash
# Check matching up/down pairs
for f in db/migrations/*.up.sql; do
  down="${f/.up.sql/.down.sql}"
  [ -f "$down" ] || echo "MISSING: $down"
done
```

### 10. Security Scan

```bash
golangci-lint run --enable gosec --no-config ./internal/... 2>&1
```

Run gosec standalone against internal packages. Flag any findings above informational level. Common things to watch for in this project:
- SQL injection via string concatenation (should use parameterized queries)
- Hardcoded credentials (false positives on `token`/`password` field names are OK)
- Weak crypto usage

## Output Format

After running all checks, present a summary:

```
## Verification Results

| # | Check               | Status | Details          |
|---|---------------------|--------|------------------|
| 1 | Build               | PASS   |                  |
| 2 | Go Vet              | PASS   |                  |
| 3 | Lint (0 warnings)   | PASS   | 18 linters       |
| 4 | Format              | PASS   |                  |
| 5 | Module Tidy         | PASS   |                  |
| 6 | Tests               | PASS   | 35 passed        |
| 7 | Race Detection      | PASS   |                  |
| 8 | Coverage            | PASS   | 74.9% overall    |
| 9 | Migration Files     | PASS   | 3 pairs          |
|10 | Security Scan       | PASS   |                  |

**Result: ALL PASSED**
```

If anything fails:

```
**Result: FAILED — fix items marked FAIL above, then run `/verify` again.**
```

## Auto-fix Behavior

For these checks, auto-fix without asking:
- **Format**: run `go fmt ./...`
- **Module tidy**: run `go mod tidy`
- **Lint**: fix issues if they are straightforward (combine params, pre-allocate slices, etc.)

For these checks, report but do NOT auto-fix:
- **Test failures**: the user needs to understand what broke
- **Race conditions**: require design-level thinking
- **Coverage gaps**: user decides whether to add tests or adjust thresholds
- **Security findings**: user must evaluate risk
