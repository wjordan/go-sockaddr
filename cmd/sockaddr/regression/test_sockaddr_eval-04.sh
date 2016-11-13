#!/bin/sh --

set -e
exec 2>&1
cat <<'EOF' | exec ../sockaddr eval -
{{GetIfSockAddrs | include "name" "lo0" | printf "%v"}}
EOF
