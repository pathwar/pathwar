Pathwar Core [![Build Status](https://travis-ci.org/pathwar/core.svg?branch=master)](https://travis-ci.org/pathwar/core)
============

## Deprecated! Go to the Pathwar monorepo: https://github.com/pathwar/pathwar

Set of tools to create levels on Pathwar.

We worked hard to make the levels runnable on any machine running docker
and we are still working hard to make the publishing as easy as possible.

This repository will help you to begin in the level-writing, don't hesitate to ask anything, or open an [issue](https://github.com/pathwar/core/issues).

---

Tools
-----

- [skeleton](https://github.com/pathwar/core/tree/master/skeleton): Contains the `/pathwar` directory
- [mk](https://github.com/pathwar/core/tree/master/mk): Contains the `.mk` files included by `Makefile` on levels and improve the confort of level development
- [templates](https://github.com/pathwar/core/tree/master/templates): Templates for common levels in multiple languages with examples

---

Templates
---------

Name        | Versions            | Links
------------|---------------------|-----------------
apache-php  | 5.6.4               | [Code](templates/apache-php), [Examples](templates/apache-php/examples)
nodejs      | 0.10.34             | [Code](templates/nodejs), [Examples](templates/nodejs/examples)
nginx       | 1.7.8               | [Code](templates/nginx), [Examples](templates/nginx/examples)
ruby        | 2.1.5               | [Code](templates/ruby), [Examples](templates/ruby/examples)
supervisord | n/a                 | [Code](templates/supervisord), [Examples](templates/supervisord/examples)
golang      | 1.4.3               | [Code](templates/golang), [Examples](templates/golang/examples)
phantomjs   | 1.9.8               | [Code](templates/phantomjs), [Examples](templates/phantomjs/examples)
python      | 2.7.9, 3.3.6, 3.4.2 | [Code](templates/python), [Examples](templates/python/examples), [Docker](https://registry.hub.docker.com/u/pathwar/python/)

All templates support standard and [onbuild](https://docs.docker.com/reference/builder/#onbuild) modes.

---

Level requirements (see [level-helloworld](https://github.com/pathwar/level-helloworld))
------------------

- have a `docker-compose.yml` defining the level, you can have multiple linked containers, multiple ports. [(example)](https://github.com/pathwar/level-helloworld/blob/master/docker-compose.yml)
- even if you can run levels without any `Dockerfile` using `docker-compose`, you need to have at least one Dockerfile inheriting from a `core` template (this repository), so we can integrate your level with our system. [(example)](https://github.com/pathwar/level-helloworld/blob/master/level.yml)
- a `level.yml` file (manifest) defining your level specs, except for our system. [(example)](https://github.com/pathwar/level-helloworld/blob/master/level.yml)
- a `scripts` directory with overrides for the [skeleton](https://github.com/pathwar/core/tree/master/skeleton/scripts). [(example)](https://github.com/pathwar/level-helloworld/tree/master/scripts)

Bonus:

- a screenshot to make the repo fancy
- make the level private or at least difficult to understand with the sources

Limitations:

- you cannot link host-volume on containers, but you can use volume link between containers
- your level needs to be buildable without cloning any other repository, however you can use some `wget` hacks from some scripts


---

Levels
------

Level            | Containers | Technos       | Open-Source | Repository
-----------------|------------|---------------|-------------|---------------------------------------------
helloworld       | 1          | nginx         | :o:         | http://github.com/pathwar/level-helloworld
pnu              | 1          | php           | :o:         | http://github.com/pathwar/pnu
captcha          | 1          | php           | :o:         | http://github.com/pathwar/captcha
calc             | 2          | python, mysql | :o:         | http://github.com/pathwar/calc
upload-hi        | 1          | php           | :o:         | http://github.com/pathwar/upload-hi
upload-kthxbie   | 1          | php           | :x:         | http://github.com/pathwar/upload-kthxbie
training-sqli    | 2          | php, mysql    | :o:         | http://github.com/pathwar/training-sqli
training-http    | 1          | n/a           | :o:         | http://github.com/pathwar/training-http
training-include | 1          | n/a           | :o:         | http://github.com/pathwar/training-include
training-brute   | 1          | n/a           | :o:         | http://github.com/pathwar/training-brute
training-tools   | 1          | n/a           | :o:         | http://github.com/pathwar/training-tools

---

© 2014-2015 Pathwar Team - [MIT License](https://github.com/pathwar/core/blob/master/LICENSE.md).
