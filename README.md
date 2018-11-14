# gostron
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
