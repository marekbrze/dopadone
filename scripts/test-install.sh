#!/bin/bash
#
# Test script for install.sh
# Tests installation logic with mock GitHub API responses
#
# Usage: ./scripts/test-install.sh [TEST_NAME]
#
# If TEST_NAME is provided, only that test will run.
# Otherwise, all tests run.
#
# Exit codes:
#   0 - All tests passed
#   1 - One or more tests failed

# shellcheck disable=SC2329,SC2016
# Above directives suppress false positives:
# - SC2329: Functions are invoked indirectly via run_test
# - SC2016: We intentionally use single quotes to grep for literal strings

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_SCRIPT="${SCRIPT_DIR}/install.sh"

PASSED=0
FAILED=0
TESTS_RUN=0
SKIPPED=0

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    PASSED=$((PASSED + 1))
}

log_fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    echo "  Error: $2"
    FAILED=$((FAILED + 1))
}

log_skip() {
    echo -e "${YELLOW}⊘ SKIP${NC}: $1"
    SKIPPED=$((SKIPPED + 1))
}

log_info() {
    echo -e "${YELLOW}ℹ INFO${NC}: $1"
}

run_test() {
    local test_name="$1"
    local test_func="$2"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    echo ""
    echo "Running: $test_name"
    echo "----------------------------------------"
    
    if $test_func; then
        log_pass "$test_name"
        return 0
    else
        return 1
    fi
}

test_help_flag() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --help 2>&1); then
        if echo "$output" | grep -q "Usage:" && \
           echo "$output" | grep -q "\-\-dry-run" && \
           echo "$output" | grep -q "\-\-yes" && \
           echo "$output" | grep -q "INSTALL_DIR"; then
            return 0
        else
            log_fail "$test_name" "Help output missing expected content"
            return 1
        fi
    else
        log_fail "$test_name" "Help flag failed"
        return 1
    fi
}

test_unknown_flag() {
    local output
    local exit_code
    
    output=$("$INSTALL_SCRIPT" --invalid-flag 2>&1) || exit_code=$?
    
    if [ "${exit_code:-0}" -eq 1 ] && echo "$output" | grep -q "Unknown option"; then
        return 0
    else
        log_fail "$test_name" "Unknown flag should exit with code 1"
        return 1
    fi
}

test_dry_run_shows_platform() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -qE "Platform:\s+(darwin|linux)-(amd64|arm64)"; then
            return 0
        else
            echo "Platform detection output not found"
            echo "Output was: $output"
            return 1
        fi
    else
        echo "Dry run failed"
        return 1
    fi
}

test_dry_run_shows_version() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -qE "Version:\s+v[0-9]+\.[0-9]+\.[0-9]+"; then
            return 0
        else
            echo "Version output not found or invalid format"
            echo "Output was: $output"
            return 1
        fi
    else
        echo "Dry run failed"
        return 1
    fi
}

test_dry_run_shows_url() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -qE "Download URL:\s+https://github\.com/.+/releases/download/v[0-9.]+/dopa-.+\.(tar\.gz|zip)"; then
            return 0
        else
            echo "Download URL not found or invalid format"
            echo "Output was: $output"
            return 1
        fi
    else
        echo "Dry run failed"
        return 1
    fi
}

test_dry_run_no_download() {
    local output
    
    output=$("$INSTALL_SCRIPT" --dry-run 2>&1)
    
    if echo "$output" | grep -q "No files were downloaded or modified"; then
        return 0
    else
        echo "Dry run should indicate no files were modified"
        echo "Output was: $output"
        return 1
    fi
}

test_platform_detection_current() {
    local output
    local expected_os
    local expected_arch
    local expected_platform
    
    expected_os=$(uname -s | tr '[:upper:]' '[:lower:]')
    expected_arch=$(uname -m)
    
    case "$expected_arch" in
        x86_64|amd64) expected_arch="amd64" ;;
        aarch64|arm64) expected_arch="arm64" ;;
    esac
    
    expected_platform="${expected_os}-${expected_arch}"
    
    if output=$("$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -q "Platform:.*${expected_platform}"; then
            return 0
        else
            echo "Expected platform: $expected_platform"
            echo "Output was: $output"
            return 1
        fi
    else
        echo "Dry run failed"
        return 1
    fi
}

test_custom_install_dir() {
    local output
    local custom_dir="/tmp/dopa-test-install"
    
    if output=$(INSTALL_DIR="$custom_dir" "$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -q "Install to:.*${custom_dir}/dopa"; then
            return 0
        else
            echo "Custom INSTALL_DIR not reflected in output"
            echo "Output was: $output"
            return 1
        fi
    else
        echo "Dry run with custom INSTALL_DIR failed"
        return 1
    fi
}

test_yes_flag_in_dry_run() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run --yes 2>&1); then
        return 0
    else
        echo "--yes flag should work in dry-run mode"
        return 1
    fi
}

test_no_verify_flag_in_dry_run() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run --no-verify 2>&1); then
        return 0
    else
        echo "--no-verify flag should work in dry-run mode"
        return 1
    fi
}

test_dependency_check_function() {
    if grep -q "check_dependencies()" "$INSTALL_SCRIPT"; then
        if grep -q "command -v curl" "$INSTALL_SCRIPT" && \
           grep -q "command -v tar" "$INSTALL_SCRIPT" && \
           grep -q "command -v unzip" "$INSTALL_SCRIPT"; then
            return 0
        else
            echo "check_dependencies() missing required dependency checks"
            return 1
        fi
    else
        echo "check_dependencies() function not found"
        return 1
    fi
}

test_verify_installation_function() {
    if grep -q "verify_installation()" "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "verify_installation() function not found"
        return 1
    fi
}

test_upgrade_support_function() {
    if grep -q "check_existing_installation()" "$INSTALL_SCRIPT" && \
       grep -q "prompt_upgrade()" "$INSTALL_SCRIPT" && \
       grep -q "backup_existing()" "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "Missing upgrade support functions"
        return 1
    fi
}

test_binary_rename_logic() {
    if grep -q 'binary_in_archive="${BINARY_NAME}-${platform}' "$INSTALL_SCRIPT" || \
       grep -q 'binary_in_archive="${BINARY_NAME}-\${platform}' "$INSTALL_SCRIPT" || \
       grep -q 'binary_in_archive=.*BINARY_NAME.*platform' "$INSTALL_SCRIPT"; then
        if grep -q 'mv.*binary_in_archive.*extracted_binary' "$INSTALL_SCRIPT" || \
           grep -q 'mv "\$binary_in_archive" "\$extracted_binary"' "$INSTALL_SCRIPT"; then
            return 0
        else
            echo "Binary rename logic not found"
            return 1
        fi
    else
        echo "Binary naming pattern not correctly handled"
        return 1
    fi
}

test_windows_binary_handling() {
    if grep -q 'BINARY_NAME.*\.exe' "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "Windows binary (.exe) handling not found"
        return 1
    fi
}

test_error_handling() {
    if grep -q "set -e" "$INSTALL_SCRIPT"; then
        if grep -q "exit 1" "$INSTALL_SCRIPT" && grep -q "exit 2" "$INSTALL_SCRIPT"; then
            return 0
        else
            echo "Missing proper exit codes"
            return 1
        fi
    else
        echo "Script doesn't use set -e for error handling"
        return 1
    fi
}

test_url_generation() {
    if grep -q "get_download_url()" "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "URL generation function not found"
        return 1
    fi
}

test_url_format() {
    local output
    
    if output=$("$INSTALL_SCRIPT" --dry-run 2>&1); then
        if echo "$output" | grep -q "releases/download/v.*dopa-.*\\.tar\\.gz"; then
            return 0
        else
            echo "URL format doesn't match expected pattern"
            echo "Output: $output"
            return 1
        fi
    else
        return 1
    fi
}

test_existing_installation_detection() {
    if grep -q "check_existing_installation" "$INSTALL_SCRIPT" && \
       grep -q "Existing installation" "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "Existing installation detection not found"
        return 1
    fi
}

test_backup_functionality() {
    if grep -q "backup_existing" "$INSTALL_SCRIPT" && \
       grep -q "\.backup\." "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "Backup functionality not found"
        return 1
    fi
}

test_confirmation_prompt() {
    if grep -q "read.*-p.*Replace" "$INSTALL_SCRIPT" || \
       grep -q "prompt_upgrade" "$INSTALL_SCRIPT"; then
        return 0
    else
        echo "Confirmation prompt for upgrades not found"
        return 1
    fi
}

ALL_TESTS="help_flag unknown_flag dry_run_platform dry_run_version dry_run_url dry_run_no_download platform_current custom_install_dir yes_flag no_verify_flag dependency_check verify_installation upgrade_support binary_rename windows_handling error_handling url_generation url_format existing_installation backup_functionality confirmation_prompt"

run_all_tests() {
    local test_name
    
    for test_name in $ALL_TESTS; do
        case "$test_name" in
            help_flag) run_test "help_flag" test_help_flag ;;
            unknown_flag) run_test "unknown_flag" test_unknown_flag ;;
            dry_run_platform) run_test "dry_run_platform" test_dry_run_shows_platform ;;
            dry_run_version) run_test "dry_run_version" test_dry_run_shows_version ;;
            dry_run_url) run_test "dry_run_url" test_dry_run_shows_url ;;
            dry_run_no_download) run_test "dry_run_no_download" test_dry_run_no_download ;;
            platform_current) run_test "platform_detection_current" test_platform_detection_current ;;
            custom_install_dir) run_test "custom_install_dir" test_custom_install_dir ;;
            yes_flag) run_test "yes_flag" test_yes_flag_in_dry_run ;;
            no_verify_flag) run_test "no_verify_flag" test_no_verify_flag_in_dry_run ;;
            dependency_check) run_test "dependency_check" test_dependency_check_function ;;
            verify_installation) run_test "verify_installation" test_verify_installation_function ;;
            upgrade_support) run_test "upgrade_support" test_upgrade_support_function ;;
            binary_rename) run_test "binary_rename" test_binary_rename_logic ;;
            windows_handling) run_test "windows_handling" test_windows_binary_handling ;;
            error_handling) run_test "error_handling" test_error_handling ;;
            url_generation) run_test "url_generation" test_url_generation ;;
            url_format) run_test "url_format" test_url_format ;;
            existing_installation) run_test "existing_installation" test_existing_installation_detection ;;
            backup_functionality) run_test "backup_functionality" test_backup_functionality ;;
            confirmation_prompt) run_test "confirmation_prompt" test_confirmation_prompt ;;
        esac
    done
}

run_single_test() {
    local test_name="$1"
    
    case "$test_name" in
        help_flag) run_test "help_flag" test_help_flag ;;
        unknown_flag) run_test "unknown_flag" test_unknown_flag ;;
        dry_run_platform) run_test "dry_run_platform" test_dry_run_shows_platform ;;
        dry_run_version) run_test "dry_run_version" test_dry_run_shows_version ;;
        dry_run_url) run_test "dry_run_url" test_dry_run_shows_url ;;
        dry_run_no_download) run_test "dry_run_no_download" test_dry_run_no_download ;;
        platform_current) run_test "platform_detection_current" test_platform_detection_current ;;
        custom_install_dir) run_test "custom_install_dir" test_custom_install_dir ;;
        yes_flag) run_test "yes_flag" test_yes_flag_in_dry_run ;;
        no_verify_flag) run_test "no_verify_flag" test_no_verify_flag_in_dry_run ;;
        dependency_check) run_test "dependency_check" test_dependency_check_function ;;
        verify_installation) run_test "verify_installation" test_verify_installation_function ;;
        upgrade_support) run_test "upgrade_support" test_upgrade_support_function ;;
        binary_rename) run_test "binary_rename" test_binary_rename_logic ;;
        windows_handling) run_test "windows_handling" test_windows_binary_handling ;;
        error_handling) run_test "error_handling" test_error_handling ;;
        url_generation) run_test "url_generation" test_url_generation ;;
        url_format) run_test "url_format" test_url_format ;;
        existing_installation) run_test "existing_installation" test_existing_installation_detection ;;
        backup_functionality) run_test "backup_functionality" test_backup_functionality ;;
        confirmation_prompt) run_test "confirmation_prompt" test_confirmation_prompt ;;
        *)
            echo "Error: Unknown test '$test_name'"
            echo "Available tests: $ALL_TESTS"
            exit 1
            ;;
    esac
}

main() {
    local specific_test="$1"
    
    echo "========================================"
    echo "  Dopadone Install Script Test Suite"
    echo "========================================"
    
    if [ ! -f "$INSTALL_SCRIPT" ]; then
        echo "Error: Install script not found at $INSTALL_SCRIPT"
        exit 1
    fi
    
    if [ -n "$specific_test" ]; then
        run_single_test "$specific_test"
    else
        run_all_tests
    fi
    
    echo ""
    echo "========================================"
    echo "  Test Results"
    echo "========================================"
    echo "Tests run: $TESTS_RUN"
    echo -e "Passed:    ${GREEN}$PASSED${NC}"
    echo -e "Failed:    ${RED}$FAILED${NC}"
    if [ "$SKIPPED" -gt 0 ]; then
        echo -e "Skipped:   ${YELLOW}$SKIPPED${NC}"
    fi
    echo ""
    
    if [ "$FAILED" -gt 0 ]; then
        exit 1
    fi
    
    exit 0
}

main "$@"
