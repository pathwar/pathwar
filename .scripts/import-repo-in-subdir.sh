#!/bin/sh -xe

REPO="$1"
SUBDIR="$2"
BRANCH="$3"
GIT_PREFIX=${GIT_PREFIX:-"git@github.com:pathwar"}
REMOTEREPO=${REMOTEREPO:-pathwar}
SUBDIR=${SUBDIR:-$REPO}
BRANCH=${BRANCH:-master}
BRANCH_CLEAN=$(echo "$BRANCH" | tr "/" "-")

rm -rf tmp-${REPO}-${BRANCH_CLEAN}
git clone "${GIT_PREFIX}/${REPO}" tmp-${REPO}-${BRANCH_CLEAN}

cd tmp-${REPO}-${BRANCH_CLEAN}
git checkout ${BRANCH}
git filter-branch --tree-filter         'mkdir -p __tmp__; mv $(ls -a1 | sed 1d | sed 1d | grep -v __tmp__) __tmp__/' HEAD
git filter-branch --force --tree-filter 'mkdir -p '$SUBDIR'; mv $(find __tmp__/ -maxdepth 1 | sed 1d) '$SUBDIR'/' HEAD

git remote add ${REMOTEREPO} ${GIT_PREFIX}/${REMOTEREPO}.git
git checkout -b sync-$REPO-${BRANCH}
git push -u ${REMOTEREPO} sync-$REPO-${BRANCH}
