#!/bin/bash

if [ $# -eq 0 ]
  then
    echo "Must provide Skeleton Directory"
    exit
fi

SKELETONDIR=${1}
CURRENTDIR=$(dirname ${0})

mkdir -p ${SKELETONDIR}
rm -rf ${SKELETONDIR}/* ${SKELETONDIR}/.??*

cp -rf ${CURRENTDIR}/scripts ${SKELETONDIR}/
