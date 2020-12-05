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
BUILD_DIR=${TOPDIR}/build
MGMT_COMMON_DIR=$(realpath ${TOPDIR}/../sonic-mgmt-common)


# Setup database config file path
if [[ -z ${DB_CONFIG_PATH} ]]; then
    F=${TOPDIR}/testdata/database_config.json
    G=${BINDIR}/database_config.json
    S=$(grep unix_socket_path $F | sed -E 's/.*:\s*"(.*)".*/\1/')

    # Reuse already generated build/bin/database_config.json if any
    if [[ -e $G ]]; then
        export DB_CONFIG_PATH=$G
    # Use testdata/database_config.json as is if its socket path is valid
    elif [[ -e $S ]]; then
        export DB_CONFIG_PATH=$F
    # Try to guess a redis socket path; and generate build/bin/database_config.json
    # by replacing the socket config in testdata/database_config.json
    else
        GUESS=( /tmp/redis.sock /tmp/redis/redis.sock /var/run/redis/redis.sock \
            /var/run/redis/redis-server.sock "${BINDIR}/redis.sock" )
        for S in "${GUESS[@]}"; do
            [[ -e $S ]] && P=${S//\//\\\/} && break
        done
        if [[ -z $P ]]; then
            >&2 echo "error: Could not guess redis socket path."
            >&2 echo "Please update '$(realpath --relative-to=. $F)' and retry."
            exit 1
        fi
        sed -E "s/\"unix_socket_path\".*/\"unix_socket_path\": \"$P\"/" $F > $G
        export DB_CONFIG_PATH=$G
    fi
fi

# LD_LIBRARY_PATH for CVL
if [[ -z ${LD_LIBRARY_PATH} ]]; then
    export LD_LIBRARY_PATH=/usr/local/lib
fi

# Setup CVL schema directory
if [[ -z ${CVL_SCHEMA_PATH} ]]; then
    export CVL_SCHEMA_PATH=${MGMT_COMMON_DIR}/build/cvl/schema
fi

# Prepare CVL config file with all traces enabled
if [[ -z ${CVL_CFG_FILE} ]]; then
    export CVL_CFG_FILE=${BINDIR}/cvl_cfg.json
    if [[ ! -e ${CVL_CFG_FILE} ]]; then
        F=${MGMT_COMMON_DIR}/cvl/conf/cvl_cfg.json
        sed -E 's/((TRACE|LOG).*)\"false\"/\1\"true\"/' $F > ${CVL_CFG_FILE}
    fi
fi

# Prepare yang files directiry for transformer
if [[ -z ${YANG_MODELS_PATH} ]]; then
    export YANG_MODELS_PATH=${BUILD_DIR}/all_yangs
    mkdir -p ${YANG_MODELS_PATH}
    pushd ${YANG_MODELS_PATH} > /dev/null
    MGMT_COMN=$(realpath --relative-to=${PWD} ${MGMT_COMMON_DIR})
    rm -f *
    find ${MGMT_COMN}/models/yang -name "*.yang" -not -path "*/testdata/*" -exec ln -sf {} \;
    ln -sf ${MGMT_COMN}/models/yang/version.xml
    ln -sf ${MGMT_COMN}/config/transformer/models_list
    ln -sf ${MGMT_COMN}/config/transformer/sonic_table_info.json
    ln -sf ${MGMT_COMN}/build/yang/api_ignore
    popd > /dev/null
fi

echo "CVL_SCHEMA_PATH  = ${CVL_SCHEMA_PATH}"
echo "CVL_CFG_FILE     = ${CVL_CFG_FILE}"
echo "DB_CONFIG_PATH   = ${DB_CONFIG_PATH}"
echo "YANG_MODELS_PATH = ${YANG_MODELS_PATH}"

