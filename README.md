# Gostron
golang webserver with pluggable route controller

## Scope
It run only on linux because of go plugin support.

## Usage
- create middleware plugins under plugins/middlewares
- create handler plugins under plugins/handlers

- configure routes in the configurations/routes.json file
- configure server in the configurations/server.json file

### Tools
the create.sh script provide scaffold to make your middlewares and handlers.
executing command:

```bash
  $ ./create.sh handler mio
```
it will produce a plugins/handlers/mio.go file with the structure needed to use it in the server, as the same of

```bash
  $ ./create.sh middleware mio
```
that will create the plugins/middlewares/mio.go file.


### Building
once you have finish configurations and created the handlers/middlewares plugins, in shell run the command:

```bash
  $ make build
```

if you want to remove all compiled files, run:

```bash
  $ make clean
```

## TODO
### Example
- create a ipfilter middleware
- create basic auth middleware
- create a "only-admin-access" middleware
### Test
- routing configuration test: searh for duplicated or wrong path, search for required plugins
- performance test
### Desired features
- server config to REDIRECT HTTP TO HTTPS
- middleware: CLIENT AUTHENTICATION
- server config to enable HTTPS: use crypto/tls package with ability to rotate TLS session ticket keys by default
- JWT API auth for javascript frontend framework like angular
- csrf token for request validation
