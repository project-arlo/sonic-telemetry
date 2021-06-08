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
echo "usage: $(basename $0) [OPTIONS] OPERATION* [-- [gnmi_set args]]"
echo ""
echo "OPTIONS:"
echo "  -host HOST          Server IP address (default 127.0.0.1)"
echo "  -port PORT          Server port (default 8080)"
echo "  -user USER:PASS     Username and password for authentication"
echo ""
echo "OPERATION: (can be repeated)"
echo "  -delete  PATH       Delete path"
echo "  -update  PATH JSON  Update path and json value"
echo "  -replace PATH JSON  Replace path and json value"
echo ""
}

TOPDIR=$(git rev-parse --show-toplevel)
BINDIR=${TOPDIR}/build/bin
gnmi_set=$(realpath --relative-to ${PWD} ${BINDIR}/gnmi_set)

if [[ ! -f ${gnmi_set} ]]; then
    echo "error: gNMI tools are not compiled"
    echo "Please run 'make telemetry' and try again"
    exit 1
fi

HOST=localhost
PORT=8080
ARGS=()

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
    -D|-delete|--delete)
        ARGS+=( -delete "$2" )
        shift 2;;
    -U|-update|--update)
        F=$(mktemp -t 'u_XXXXX.json')
        echo "$3" > $F
        ARGS+=( -update "/${2#/}:@$F" )
        shift 3;;
    -R|-replace|--replace)
        F=$(mktemp -t 'r_XXXXX.json')
        echo "$3" > $F
        ARGS+=( -replace "/${2#/}:@$F" )
        shift 3;;
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
${gnmi_set} "${ARGS[@]}"
