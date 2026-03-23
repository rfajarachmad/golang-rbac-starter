#!/usr/bin/env bash
set -euo pipefail

# ── Coverage thresholds (%) ──────────────────────────────────────────────
# High-ROI: enforce on layers where bugs actually hide.
# Skip config (infra), entity/model (structs), repository (thin GORM wrappers).
OVERALL_MIN=70
USECASE_MIN=60
DELIVERY_MIN=70

# ── Colors ───────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

# ── Run tests with coverage ─────────────────────────────────────────────
COVER_FILE="coverage.out"
echo "Running tests with coverage..."
go test -count=1 -coverprofile="$COVER_FILE" -coverpkg=./internal/... ./test/ 2>&1 | tail -1

# ── Extract coverage percentage for a package pattern ────────────────────
get_coverage() {
    local pattern="$1"
    go tool cover -func="$COVER_FILE" \
        | grep "$pattern" \
        | awk '{print $NF}' \
        | sed 's/%//' \
        | awk '{ total += $1; count++ } END { if (count > 0) printf "%.1f", total/count; else print "0.0" }'
}

# ── Get overall coverage ────────────────────────────────────────────────
OVERALL=$(go tool cover -func="$COVER_FILE" | grep "^total:" | awk '{print $NF}' | sed 's/%//')

# ── Get per-layer coverage ──────────────────────────────────────────────
USECASE=$(get_coverage "internal/usecase/")
DELIVERY=$(get_coverage "internal/delivery/")

# ── Check thresholds ────────────────────────────────────────────────────
FAILED=0

check_threshold() {
    local name="$1"
    local actual="$2"
    local minimum="$3"

    if (( $(echo "$actual < $minimum" | bc -l) )); then
        echo -e "  ${RED}FAIL${NC}  $name: ${actual}% < ${minimum}% minimum"
        FAILED=1
    else
        echo -e "  ${GREEN}PASS${NC}  $name: ${actual}% (≥ ${minimum}%)"
    fi
}

echo ""
echo "Coverage thresholds:"
check_threshold "Overall " "$OVERALL"  "$OVERALL_MIN"
check_threshold "UseCase " "$USECASE"  "$USECASE_MIN"
check_threshold "Delivery" "$DELIVERY" "$DELIVERY_MIN"
echo ""

if [ "$FAILED" -eq 1 ]; then
    echo -e "${RED}Coverage check failed.${NC} Improve tests in the packages above."
    exit 1
else
    echo -e "${GREEN}All coverage thresholds passed.${NC}"
fi
