#!/usr/bin/env bash
################################################################################
#                                                                              #
#  Copyright 2021 Broadcom. The term Broadcom refers to Broadcom Inc. and/or   #
#  its subsidiaries.                                                           #
#                                                                              #
#  Licensed under the Apache License, Version 2.0 (the "License");             #
#  you may not use this file except in compliance with the License.            #
#  You may obtain a copy of the License at                                     #
#                                                                              #
#     http://www.apache.org/licenses/LICENSE-2.0                               #
#                                                                              #
#  Unless required by applicable law or agreed to in writing, software         #
#  distributed under the License is distributed on an "AS IS" BASIS,           #
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.    #
#  See the License for the specific language governing permissions and         #
#  limitations under the License.                                              #
#                                                                              #
################################################################################

set -e

function print_usage() {
echo "usage: $(basename $0) [OPTIONS] PATH* [-- [gnmi_get args]]"
echo ""
echo "OPTIONS:"
echo "  -host HOST        Server IP address (default 127.0.0.1)"
echo "  -port PORT        Server port (default 8080)"
echo "  -proto            Use PROTO encoding"
echo "  -user USER:PASS   Username and password for authentication"
echo ""
}

TOPDIR=$(git rev-parse --show-toplevel)
BINDIR=${TOPDIR}/build/bin
gnmi_get=$(realpath --relative-to ${PWD} ${BINDIR}/gnmi_get)

if [[ ! -f ${gnmi_get} ]]; then
    echo "error: gNMI tools are not compiled"
    echo "Please run 'make telemetry' and try again"
    exit 1
fi

HOST=localhost
PORT=8080
ARGS=()
PATHS=()

while [[ $# -gt 0 ]]; do
    case "$1" in
    -h|-help|--help)
        print_usage
        exit 0;;
    -H|-host|--host)
        HOST=$2
        shift 2;;
    -p|-port|--port)
        PORT=$2
        shift 2;;
    -u|-user|--user)
        ARGS+=( -username "${2%%:*}" -password "${2#*:}" )
        shift 2;;
    -proto|--proto)
        ARGS+=( -encoding PROTO )
        shift;;
    [/_a-zA-Z]*)
        PATHS+=( -xpath "$1" )
        shift;;
    --)
        shift
        ARGS+=( "$@" )
        break;;
    *)
        echo "error: unknown option: $1"
        print_usage
        exit 1;;
    esac
done

ARGS+=( -insecure )
[[ "$@" =~ -(also)?log*  ]] || ARGS+=( -logtostderr )
[[ "$@" =~ -target_addr* ]] || ARGS+=( -target_addr ${HOST}:${PORT} )

set -x
${gnmi_get} "${ARGS[@]}" "${PATHS[@]}"
