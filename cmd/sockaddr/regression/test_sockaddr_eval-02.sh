#!/bin/sh --

set -e
exec 2>&1
exec ../sockaddr eval '{{GetIfSockAddrs | include "name" "lo0" | printf "%v"}}' '{{GetIfSockAddrs | include "name" "lo0" | printf "%v"}}'
