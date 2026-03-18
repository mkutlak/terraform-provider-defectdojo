#!/usr/bin/env bash
set -euo pipefail

# Multi-version DefectDojo compatibility checker.
# Spins up each DD version, fetches its OpenAPI spec, and optionally runs acceptance tests.
#
# Usage:
#   ./scripts/dd-version-compat.sh [--test] [VERSION...]
#
# Examples:
#   ./scripts/dd-version-compat.sh                    # Check default versions, spec only
#   ./scripts/dd-version-compat.sh --test              # Check + run acceptance tests
#   ./scripts/dd-version-compat.sh 2.54.3 2.42.0      # Check specific versions
#   ./scripts/dd-version-compat.sh --test 2.54.3       # Test a specific version

DEFAULT_VERSIONS=("2.56.2" "2.54.3")
RUN_TESTS=false
VERSIONS=()
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Parse arguments
for arg in "$@"; do
    if [[ "$arg" == "--test" ]]; then
        RUN_TESTS=true
    else
        VERSIONS+=("$arg")
    fi
done

if [[ ${#VERSIONS[@]} -eq 0 ]]; then
    VERSIONS=("${DEFAULT_VERSIONS[@]}")
fi

# Results tracking
declare -A SPEC_RESULTS
declare -A TEST_RESULTS

cleanup() {
    echo "Cleaning up..."
    cd "$PROJECT_DIR"
    docker compose down -v 2>/dev/null || true
}
trap cleanup EXIT

for version in "${VERSIONS[@]}"; do
    echo ""
    echo "============================================"
    echo "  DefectDojo $version"
    echo "============================================"

    # Clean slate — different DD versions have incompatible DB schemas
    cd "$PROJECT_DIR"
    echo "Stopping any running containers..."
    docker compose down -v 2>/dev/null || true

    # Start this version
    echo "Starting DefectDojo $version..."
    export DD_VERSION="$version"
    docker compose up -d

    # Wait for DD to be ready
    echo "Waiting for DefectDojo API to be ready..."
    ready=false
    for i in $(seq 1 60); do
        if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v2/ 2>/dev/null | grep -qE "200|401|403"; then
            echo "DefectDojo $version is ready."
            ready=true
            break
        fi
        echo "  Attempt $i/60 - waiting 10s..."
        sleep 10
    done

    if [[ "$ready" != "true" ]]; then
        echo "FAIL: DefectDojo $version did not start in time."
        SPEC_RESULTS[$version]="FAIL (timeout)"
        TEST_RESULTS[$version]="SKIP"
        docker compose logs uwsgi nginx 2>/dev/null | tail -50
        continue
    fi

    # Fetch OpenAPI spec
    echo "Fetching OpenAPI spec..."
    if make dd-spec DD_VERSION="$version" 2>/dev/null; then
        SPEC_RESULTS[$version]="OK"
        echo "Spec saved to openapi-specs/$version/defect_dojo.json"
    else
        SPEC_RESULTS[$version]="FAIL"
        echo "FAIL: Could not fetch OpenAPI spec for $version"
    fi

    # Optionally run acceptance tests
    if [[ "$RUN_TESTS" == "true" ]]; then
        echo "Running acceptance tests against DefectDojo $version..."
        if make testacc-local 2>&1 | tee "/tmp/dd-compat-test-$version.log"; then
            TEST_RESULTS[$version]="PASS"
        else
            TEST_RESULTS[$version]="FAIL"
            echo "Test log saved to /tmp/dd-compat-test-$version.log"
        fi
    else
        TEST_RESULTS[$version]="SKIP"
    fi

    # Stop before next version
    echo "Stopping DefectDojo $version..."
    docker compose down -v
done

# Print summary
echo ""
echo "============================================"
echo "  Compatibility Results"
echo "============================================"
printf "%-20s %-15s %-15s\n" "VERSION" "SPEC" "TESTS"
printf "%-20s %-15s %-15s\n" "-------" "----" "-----"
for version in "${VERSIONS[@]}"; do
    printf "%-20s %-15s %-15s\n" "$version" "${SPEC_RESULTS[$version]:-N/A}" "${TEST_RESULTS[$version]:-N/A}"
done
echo ""
