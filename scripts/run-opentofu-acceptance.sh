#!/usr/bin/env bash
#
# Run provider acceptance tests against the OpenTofu CLI instead of Terraform.
#
# A test selection is ALWAYS required. Pass a -run pattern (with wildcards /
# regex) and, optionally, the acceptance package(s) to search. Without a
# -run pattern the script refuses to run (it will not silently default to a
# specific test).
#
# Usage:
#   ./scripts/run-opentofu-acceptance.sh -run 'TestAccVpcIPAddressGroupV3'
#   ./scripts/run-opentofu-acceptance.sh -run 'TestAccVpc.*'            # wildcard
#   ./scripts/run-opentofu-acceptance.sh -run 'TestAccEcs|TestAccVpc'   # several
#   ./scripts/run-opentofu-acceptance.sh -run 'TestAcc.*' ./opentelekomcloud/acceptance/ecs/...
#   ./scripts/run-opentofu-acceptance.sh -run 'TestAccFoo' -v -count 1
#
# The first argument must be -run <pattern>. Everything after the package
# paths is passed through to `go test` verbatim (standard go test flags).
#
# Requirements:
#   - the OpenTofu binary (default: ~/.local/bin/tofu, override with TOFU_PATH)
#   - provider test credentials, exported from test.env (OS_CLOUD, etc.)
#
# How the CLI is picked: the test harness drives OpenTofu through
# TF_ACC_TERRAFORM_PATH and re-attaches the in-process provider
# (TF_REATTACH_PROVIDERS), so no provider download or dev_overrides are
# needed. OpenTofu-specific env (from opentofu.org docs):
#   TF_ACC_PROVIDER_NAMESPACE=hashicorp
#   TF_ACC_PROVIDER_HOST=registry.opentofu.org
#
# Verbosity:
#   Quiet by default (TF_LOG_SDK_HELPER_RESOURCE=WARN): you only see go test's
#   own output (=== RUN / --- PASS / --- FAIL) plus, when a step fails, the
#   harness WARN line naming the failing tofu command (init/plan/apply/...)
#   together with that command's output.
#   Set VERBOSITY=TRACE for the full per-command trace while debugging.
set -euo pipefail

cd "$(dirname "$0")/.."

usage() {
  cat >&2 <<'EOF'
Error: a test selection is required.

Usage:
  run-opentofu-acceptance.sh -run <pattern> [packages...] [go test flags...]

  <pattern>   -run regex, supports wildcards. e.g. 'TestAccVpc.*',
              'TestAccEcs|TestAccVpc'. Required — no default test.
  [packages]  acceptance packages to search (default: ./opentelekomcloud/acceptance/...).

Examples:
  run-opentofu-acceptance.sh -run 'TestAccVpcIPAddressGroupV3'
  run-opentofu-acceptance.sh -run 'TestAccVpc.*'
  run-opentofu-acceptance.sh -run 'TestAcc.*' ./opentelekomcloud/acceptance/ecs/...
EOF
}

# --- require an explicit -run pattern ---------------------------------------
if [[ $# -eq 0 || "$1" != "-run" || -z "${2:-}" ]]; then
  usage
  exit 2
fi
RUN_PATTERN="$2"
shift 2   # drop `-run <pattern>`

# --- default package scope: the whole acceptance tree ------------------------
# Insert the default package target only if no package path was given. We
# detect a package path as the first remaining arg starting with './' or '/'.
has_package=0
for arg in "$@"; do
  if [[ "$arg" == ./* || "$arg" == /* ]]; then
    has_package=1
    break
  fi
done
if [[ "$has_package" -eq 0 ]]; then
  set -- ./opentelekomcloud/acceptance/... "$@"
fi

# --- OpenTofu-specific variables ---------------------------------------------
export TF_ACC_TERRAFORM_PATH="${TOFU_PATH:-$HOME/.local/bin/tofu}"
export TF_ACC_PROVIDER_NAMESPACE="hashicorp"
export TF_ACC_PROVIDER_HOST="registry.opentofu.org"
export TF_ACC=1

# --- provider auth (from repo test env) --------------------------------------
# shellcheck disable=SC1091
set -a
source test.env
set +a

# PreCheck requires these; keep values from test.env
export OS_REGION_NAME
export OS_AVAILABILITY_ZONE
export OS_SUBNET_NAME

# --- verbosity ----------------------------------------------------------------
# WARN (default): harness is silent on success; on failure it logs the
# "Error running Terraform CLI <command>" WARN with the command's output.
# TRACE: full per-command trace (every init/plan/apply with path + workdir).
export TF_LOG_SDK_HELPER_RESOURCE="${VERBOSITY:-WARN}"
export TF_LOG="${TF_LOG:-WARN}"
# Silence the in-process provider's "Previously configured provider being
# re-configured" chatter (harmless during tests). Real provider errors still
# surface because they log at ERROR, above this floor.
export TF_LOG_SDK_HELPER_SCHEMA="${TF_LOG_SDK_HELPER_SCHEMA:-ERROR}"

TOFU_BIN="$(command -v "$TF_ACC_TERRAFORM_PATH" || true)"
if [[ -z "$TOFU_BIN" ]]; then
  echo "ERROR: OpenTofu binary not found at $TF_ACC_TERRAFORM_PATH" >&2
  exit 1
fi
echo "Command: go test -count 1 -v -timeout 15m -run $RUN_PATTERN $*"
echo "Tofu CLI: $TOFU_BIN ($("$TOFU_BIN" version | head -1))"
echo "Log levels: HARNESS=$TF_LOG_SDK_HELPER_RESOURCE Tofu=$TF_LOG (override with VERBOSITY=TRACE for full per-command trace)"

# --- run ----------------------------------------------------------------------
# Re-prepend -run <pattern>; packages and any other go test flags follow.
go test -count 1 -v -timeout 15m -run "$RUN_PATTERN" "$@"
