#!/usr/bin/env bash

printf "\nStarting acceptance tests...\n"

# First of all, locate packages required for the run

base_branch="devel"

echo "Checking diff with \`$base_branch\`"
total_grep="^diff --git a/opentelekomcloud/common"
services_grep="^diff --git a/opentelekomcloud/services"

diffs=$(git diff "${base_branch}")

base_test_path=./opentelekomcloud/acceptance

# Rerun all tests if we need that
if echo "${diffs}" | grep -q "${total_grep}"; then
  echo "All acceptance tests will be run"
  TF_ACC=1 go test "${base_test_path}/..."
  exit $?
fi

# Check in there are any changes in services
if ! echo "${diffs}" | grep -q "${services_grep}"; then
  echo "No need for acceptance tests"
  exit 0
fi

# find services which needs running acceptance tests
all_services=$(echo "${diffs}" | grep "${services_grep}" | grep -Po '(?<=services\/)\w+(?=\/)' | uniq)
printf "The following services needs testing: \n%s\n" "${all_services}"

all_modules=$(echo "${all_services}" | sed -r "s|(.*)|${base_test_path}\/\1\/...|" | xargs go list)

TF_ACC=1 go test "${all_modules}" -count 1 -v -timeout 720m
