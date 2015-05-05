# Node-Pathwar: CLI + node client

[![Travis](https://img.shields.io/travis/scaleway/image-ubuntu.svg)](https://travis-ci.org/scaleway/image-ubuntu)
[![Dependency Status](https://img.shields.io/david/pathwar/node-pathwar.svg)](https://david-dm.org/pathwar/node-pathwar)

[![NPM Badge](https://nodei.co/npm/pathwar.png)](https://npmjs.org/package/pathwar)


Interact with Pathwar API from the command line.


## Usage

```console
$ pathwar -h

  Usage: pathwar [options] [command]


  Commands:

    version                                     show the version information
    ls|select [options] <type> [conditions...]  list objects
    cat|show [options] <item>                   show object
    rm|delete <item>                            remove an item
    touch|add <type> [fields...]                create an item
    update|patch <item> <fields...>             update an item

  Options:

    -h, --help            output usage information
    -V, --version         output the version number
    --api-endpoint <url>  set the API endpoint
    --token <token>       set the token
```

## Examples

List sessions

```console
$ pathwar ls sessions
 _created     _id        _updated     active   allow_new_organizations   allow_update_organizations   anonymous   email_domain    name          public
 5 days ago   62a585a9   3 days ago                                      true                                     *@pathwar.net   Beta
 5 days ago   7e5504b0   5 days ago   true     true                      true                         true        *@epitech.eu    Epitech2015   true
 5 days ago   dea6a8be   5 days ago   true     true                      true                                                     World         true
```


Create a new user-token

```console
$ pathwar add user-tokens
1abdf417-ad59-498a-a0d7-xxxxxxxxxxxx
```


Show its content

```console
$ pathwar cat 1abdf417-ad59-498a-a0d7-xxxxxxxxxxxx
{
  "_updated": "Tue, 05 May 2015 12:47:52 GMT",
  "scopes": "*",
  "description": "",
  "is_session": false,
  "token": "xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
  "expiry_date": "Wed, 06 May 2015 00:47:52 GMT",
  "is_admin": false,
  "user": "0d60edb5-82d2-4906-b879-04fca4c02f09",
  "_created": "Tue, 05 May 2015 12:47:52 GMT",
  "_id": "1abdf417-ad59-498a-a0d7-xxxxxxxxxxxx"
}
```


## Debug

`pathwar` uses the [debug](https://www.npmjs.com/package/debug) package.

To enable debug you can use the environment variable `DEBUG=` as :

- `DEBUG='*' pathwar ...` to see debug for `pathwar` and all dependencies
- `DEBUG='pathwar:*' scw ...` to see debug for `pathwar`

```console
$ DEBUG='*' pathwar ls sessions
  pathwar:lib GET https://api.pathwar.net/sessions?max_results=50&page=1&sort=-_updated& +0ms { method: 'GET',
  headers:
   { Accept: 'application/json',
     Authorization: 'Basic XXXXXXXXXXXXXXXXXXXXXX' },
  resolveWithFullResponse: true,
  url: 'https://api.pathwar.net/sessions?max_results=50&page=1&sort=-_updated&',
  json: true }
 _created     _id        _updated     active   allow_new_organizations   allow_update_organizations   anonymous   email_domain    name          public
 5 days ago   62a585a9   3 days ago                                      true                                     *@pathwar.net   Beta
 5 days ago   7e5504b0   5 days ago   true     true                      true                         true        *@epitech.eu    Epitech2015   true
 5 days ago   dea6a8be   5 days ago   true     true                      true                                                     World         true
```


## Install

1. Install `Node.js` and `npm` (https://nodejs.org/download/)
2. Install `pathwar`: `$ npm install -g pathwar`
3. Setup token: `$ echo token=XXXXX > ~/.pathwarrc`
4. Profit... `$ pathwar ls levels`


## License

[MIT](https://github.com/pathwar/node-pathwar/blob/master/LICENSE)
