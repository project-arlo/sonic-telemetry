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

if [[ ! -f ${BINDIR}/telemetry ]]; then
    >&2 echo "error: Telemetry server not compiled"
    >&2 echo "Please run 'make telemetry' and try again"
    exit 1
fi

source ${TOPDIR}/tools/test/env.sh ""

for V in "$@"; do
    case "$V" in
    -v|--v|-v=*|--v=*) HAS_V=1 ;;
    -server_crt|--server_crt|-server_crt=*|--server_crt=*) HAS_CERT=1 ;;
    -server_key|--server_key|-server_key=*|--server_key=*) HAS_CERT=1 ;;
    -client_auth|--client_auth|-client_auth=*|--client_auth=*) HAS_AUTH=1 ;;
    esac
done

ARGS=( )
[[ -z ${HAS_V} ]]    && ARGS+=( -v 2 )
[[ -z ${HAS_CERT} ]] && ARGS+=( -insecure )
[[ -z ${HAS_AUTH} ]] && ARGS+=( -client_auth none )
[[ "$@" =~ -(also)?log* ]] || ARGS+=( -logtostderr )

set -x
${BINDIR}/telemetry "${ARGS[@]}" "$@"
