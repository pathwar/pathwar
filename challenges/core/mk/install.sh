#!/bin/bash

BRANCH=${BRANCH:-master}

wget https://raw.githubusercontent.com/pathwar/core/${BRANCH}/mk/pathwar-level.mk
grep pathwar-level.mk .gitignore 2>/dev/null || echo pathwar-level.mk >> .gitignore
