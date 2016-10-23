#!/bin/sh -e -u --

set -e
set -u

if [ $# -ne 1 ]; then
    printf "Usage: %s [ test script ]\n\n" "$(basename $0)"
    printf "ERROR: Need a single test script to execute\n"
    exit 1
fi

test_script="${1}"
test_name="$(basename ${1} .sh)"
test_out="${test_name}.out"
expected_out="${test_name}.expected"

if [ ! -r "${test_script}" ]; then
    printf "ERROR: Test script %s does not exist\n" "${test_script}"
    exit 1
fi

if [ ! -r "${expected_out}" ]; then
    printf "ERROR: Expected test output does not exist\n" "${expected_out}"
    exit 1
fi

set +e
"./${test_script}" > "${test_out}" 2>&1

cmp -s "${expected_out}" "${test_out}"
result=$?
set -e

if [ "${result}" -eq 0 ]; then
    exit 0
fi

diff_out="${test_name}.diff"
diff -u "${test_out}" "${expected_out}" > "${diff_out}"
exit 1
