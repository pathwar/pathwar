#!/usr/bin/env python

# note: this script *must* not have other deps than docker-compose and
# standard packages, this allow us to use it everywhere without trouble.

import argparse
import os

from compose import config

class PackageLevel(object):
    """I package a Pathwar level

    - build the level
    - converts the docker-compose.yml build file to a production file
    - export the level
    - put everything in a tarball
    """
    def __init__(self):
        pass
    
    def go(self):
        self.build()
        self.prepare()
        self.export()
        self.compress()

    def build(self):
        pass

    def prepare(self):
        pass

    def export(self):
        pass

    def compress(self):
        pass
    

def main():
    parser = argparse.ArgumentParser(description='Package a Pathwar level')
    args = parser.parse_args()
    pkg = PackageLevel()
    pkg.go()

if __name__ == '__main__':
    main()
