# scops

Regularly reads a database and then send [Ara External Model](https://github.com/af83/ara-external-models) structures with [Protocol Buffers](https://developers.google.com/protocol-buffers/)

## :warning: Project reallocation

Scops is now located here: https://bitbucket.org/enroute-mobi/scops

## Run

To run, this app needs a plugin.

```scops -plugin path/to/plugin.so```

A make command exists to fetch a plugin from the enroute-mobi bitbucket repositories:

```make run PLUGIN=example```

### Plugin

The plugin must have a `Feeder` object that respond to this interface:

```
type Feeder interface {
	DbConnect() *dbr.Session
	GetCompleteModel(sess *dbr.Session) (*external_models.ExternalCompleteModel, error)
}
```

### Dependency

Scops and the plugings database access needs to use [the dbr library](https://github.com/gocraft/dbr).

### Options

All options can be either specified with an environment variable (SCOPS_VARIABLE=value) or via command line (-variable value).

**debug**: Boolean, enable debug messages

**syslog**: Boolean, redirect messages to syslog

**gzip**: Boolean, gzip requests

**remote**: Remote URL to send messages to

**token**: Authorization token (Authorization: Token token=secret in the request header)

**plugin**: Path to the plugin to use to get the data

**cycle**: Cycle duration (accepts a value acceptable to time.ParseDuration, ex: '60s')

