#!/bin/sh -e --

set -e
exec 2>&1
exec ../sockaddr eval '{{GetIfSockAddrs | includeByIfName "lo0" | printf "%v"}}'
