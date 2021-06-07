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
GO=${GO:-go}

TESTPKG=github.com/Azure/sonic-telemetry/gnmi_server
TESTBIN=build/tests/gnmi_server/server.test
TESTDIR=${TOPDIR}/$(dirname ${TESTBIN})
TESTARGS=()

while [[ $# -gt 0 ]]; do
    case "$1" in
    -h|-help|--help)
        echo "usage: $(basename $0) [-json] [-run NAME]"
        exit 0 ;;
    -json) JSON=1; SFLAG="-s"; shift ;;
    -run)  TESTCASE="$2"; shift 2 ;;
    *)     TESTARGS+=("$1"); shift;;
    esac
done

# Build test binary if required
cd ${TOPDIR}
make -sq ${TESTBIN} || make ${SFLAG} ${TESTBIN}

# Prepare test command
CMD=( ./$(basename ${TESTBIN}) -test.v )
[[ -z ${JSON} ]] || CMD=( ${GO} tool test2json -p "${TESTPKG}" -t "${CMD[@]}" )
[[ -z ${TESTCASE} ]] || CMD+=( -test.run "${TESTCASE}" )
[[ "${TESTARGS[@]}" =~ -log* ]] || CMD+=( -log_dir log )

mkdir -p "${TESTDIR}/log"
source ${TOPDIR}/tools/test/env.sh ${SFLAG}

[[ -z ${SFLAG} ]] && set -x
cd "${TESTDIR}"
"${CMD[@]}" "${TESTARGS[@]}"

