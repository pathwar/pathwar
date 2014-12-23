#!/bin/bash

SKELETONDIR=${1}
CURRENTDIR=$(dirname ${0})

mkdir -p ${SKELETONDIR}
rm -rf ${SKELETONDIR}/* ${SKELETONDIR}/.??*

cp -rf ${CURRENTDIR}/scripts ${SKELETONDIR}/
