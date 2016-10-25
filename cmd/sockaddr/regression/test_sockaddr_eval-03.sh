#!/bin/sh --

set -e
exec 2>&1
exec ../sockaddr eval '{{. | includeByIfName "lo0" | includeByType "IPv6" | ifAddrs | sortByAddr | joinAddrs " "}}'
