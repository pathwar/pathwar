#!/bin/bash

BRANCH=${BRANCH:-master}
SCRIPTSDIR=${SCRIPTSDIR:-/pathwar/scripts}

apply_flavor() {
    flavor="${1}"
    tar --strip=2 -C ${SCRIPTSDIR} -xzvf <(wget --no-check-certificate -qO - https://github.com/pathwar/level-scripts/archive/${BRANCH}.tar.gz) level-scripts-${BRANCH}/scripts${flavor}
}

# scripts
apply_flavor ""

# handle $FLAVORS env var
for flavor in ${FLAVORS}; do apply_flavor -${flavor}; done
