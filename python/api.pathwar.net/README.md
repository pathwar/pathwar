api.pathwar.net [![Build Status](https://travis-ci.org/pathwar/api.pathwar.net.svg?branch=master)](https://travis-ci.org/pathwar/api.pathwar.net)
===============

Pathwar official API

---

Official portal
---------------

The official [portal](https://github.com/pathwar/portal.pathwar.net/) is a consumer of this api.

There is a developement fig on [gist](https://gist.github.com/moul/fd478020ba24313359b3).

---

The official API is closed source, this repository contains resources to use the API

- [Documentation on apiary.io](http://docs.pathwar.apiary.io)
- [HTML rendering of the blueprint documentation](http://pathwar.github.io/api.pathwar.net/)

---

Resources
---------

Path                  | Resource                | Methods
----------------------|-------------------------|---------
/achievements         | Achievement collection  | GET
/achievements/{id}    | Achievement             | GET
/coupons              | Coupon collection       | GET
/coupons/{id}         | Coupon                  | GET
/levels               | Level collection        | GET
/levels/{id}          | Level                   | GET
/organizations        | Organization collection | GET
/organizations/{id}   | Organization            | GET
/sessions             | Session collection      | GET
/sessions/{id}        | Session                 | GET
/users                | User collection         | GET
/users/{id}           | User                    | GET

---

Methods
-------

Path                  | Method | Action
----------------------|--------|--------------------------
/achievements         | GET    | List achievements
/achievements/{id}    | GET    | Retrieve an achievement
/coupons              | GET    | List coupons
/coupons/{id}         | GET    | Retrieve a coupon
/levels               | GET    | List levels
/levels/{id}          | GET    | Retrieve a level
/organizations        | GET    | List organizations
/organizations/{id}   | GET    | Retrieve an organization
/sessions             | GET    | List sessions
/sessions/{id}        | GET    | Retrieve a session
/users                | GET    | List users
/users/{id}           | GET    | Retrieve a user


---

© 2014-2015 Pathwar Team - [MIT License](https://github.com/pathwar/api.pathwar.net/blob/master/LICENSE.md).
