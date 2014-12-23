#!/bin/bash

BRANCH=${BRANCH:-master}
SCRIPTSDIR=${SCRIPTSDIR:-/pathwar/scripts}
DL=${DL:-wget}

mkdir -p $SCRIPTSDIR

dl_wget() {
    wget -O - --no-check-certificate $@
}

dl_curl() {
    curl -Lk $@
}

dl() {
    dl_$DL $@
}

apply_flavor() {
    flavor="${1}"
    tar --strip=3 -C ${SCRIPTSDIR} -xzvf <(dl https://github.com/pathwar/level-templates/archive/${BRANCH}.tar.gz) level-scripts-${BRANCH}/skeleton/scripts${flavor}
}

# scripts
apply_flavor ""

# handle $FLAVORS env var
for flavor in ${FLAVORS}; do apply_flavor -${flavor}; done
