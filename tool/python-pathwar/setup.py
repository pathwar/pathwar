#!/usr/bin/env python
# -*- coding: utf-8 -*-

import imp
import os
import re

from setuptools import setup, find_packages


MODULE_NAME = 'pathwar'


def get_version():
    with open(os.path.join(
            os.path.dirname(__file__), MODULE_NAME, '__init__.py')
    ) as init:
        for line in init.readlines():
            res = re.match(r'__version__ *= *[\'"]([0-9\.]*)[\'"]$', line)
            if res:
                return res.group(1)


def get_long_description():
    readme = os.path.join(os.path.dirname(__file__), 'README.md')
    return open(readme).read()


setup(
    name=MODULE_NAME,
    version=get_version(),
    description='Pathwar API client',
    long_description=get_long_description(),
    author='Pathwar team',
    author_email='team@pathwar.net',
    license='MIT',

    install_requires=[
        'slumber >= 0.6.0',
    ],

    packages=find_packages(),

    # tests_require
    # test_suite

    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'Operating System :: POSIX',
        'Operating System :: MacOS',
        'Operating System :: Unix',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python',
        'Topic :: Software Development :: Libraries :: Python Modules',
    ],

    enty_points={
        'console_scripts': [
        ]
    }
)
