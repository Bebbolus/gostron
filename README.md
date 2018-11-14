# Gostron
golang webserver with pluggable route controller

## Scope
It run only on linux because of go plugin support.

## Usage
- create middleware plugins under plugins/middlewares
- create handler plugins under plugins/handlers

- configure routes in the configurations/routes.json file
- configure server in the configurations/server.json file

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
- https server
### Scaffold creation
- create shell command that make a pre-compiled file for handler
- create shell command that make a pre-compiled file for middleware
### Test
- server configuration test
- plugin test: if the package don't export the right function and variables, fail
- routing configuration test: searh for duplicated or wrong path, search for required plugins
- http test
- performance test
