#!/bin/sh --

set -e
set -u

for f in $(find . -name 'test_*.sh'); do
    set +e
    ./run_one.sh $f
    set -e
done

diffs=$(find . -name '*.diff')
if [ -z "${diffs}" ]; then
    exit 0
fi

printf "The following tests failed:\n\n"
for d in ${diffs}; do
    printf "\t%s\n" "$(basename ${d} .diff)"
done
exit 1
