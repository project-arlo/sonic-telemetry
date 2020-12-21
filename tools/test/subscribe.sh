#!/usr/bin/env bash

################################################################################
#                                                                              #
#  Copyright 2020 Broadcom. The term Broadcom refers to Broadcom Inc. and/or   #
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

TOPDIR=$(git rev-parse --show-toplevel || echo ${PWD})
BINDIR=${TOPDIR}/build/bin
GNMCLI=$(realpath --relative-to ${PWD} ${BINDIR}/gnmi_cli)

if [[ ! -f ${GNMCLI} ]]; then
    >&2 echo "error: gNMI tools were not compiled"
    >&2 echo "Please run 'make telemetry' and try again"
    exit 1
fi

HOST=localhost
PORT=8080
ARGS=()
PATHS=()

while [[ $# -gt 0 ]]; do
    case "$1" in
    -h|-help|--help)
        echo "usage: $(basename $0) [-H HOST] [-p PORT] [-pass] [-once|-onchange|-poll SECS|-sample SECS] PATH*"
        echo ""
        exit 0;;
    -once)
        ARGS+=( -query_type once )
        shift;;
    -onchange|-on-change|-on_change)
        ARGS+=( -query_type streaming )
        ARGS+=( -streaming_type ON_CHANGE )
        shift;;
    -sample)
        ARGS+=( -query_type streaming )
        ARGS+=( -streaming_type SAMPLE )
        ARGS+=( -streaming_sample_interval $2 )
        shift 2;;
    -poll)
        ARGS+=( -query_type polling )
        ARGS+=( -polling_interval $2s )
        shift 2;;
    -pass)
        ARGS+=( -with_user_pass )
        shift;;
    -H|-host) HOST=$2; shift 2;;
    -p|-port) PORT=$2; shift 2;;
    *) PATHS+=("$1"); shift;;
    esac
done

if [[ -z ${PATHS} ]]; then
    echo "error: At least one path required"
    exit 1
fi

ARGS+=( -insecure )
ARGS+=( -logtostderr )
ARGS+=( -address ${HOST}:${PORT} )
ARGS+=( -target OC_YANG )
ARGS+=( -query $(IFS=,; echo "${PATHS[*]}") )

set -x
${GNMCLI} "${ARGS[@]}"