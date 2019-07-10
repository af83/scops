# scops

Regularly reads a database and then send [Ara External Model](https://github.com/af83/ara-external-models) structures with [Protocol Buffers](https://developers.google.com/protocol-buffers/)

## Run

To run, this app needs a plugin.

```make run PLUGIN=example```

### Plugin

The plugin must have a `Feeder` object that respond to this interface:

```
type Feeder interface {
	DbConnect() *dbr.Session
	GetCompleteModel(sess *dbr.Session) (*external_models.ExternalCompleteModel, error)
}
```

### Options

All options can be either specified with an environment variable (SCOPS_VARIABLE=value) or via command line (-variable value).

**Debug**: Boolean, enable debug messages

**Syslog**: Boolean, redirect messages to syslog

**RemoteUrl**: Remote URL to send messages to

**AuthToken**: Authorization token (Authorization: Token token=secret in the request header)

**Plugin**: Path to the plugin to use to get the data

**Cycle**: Cycle duration (accepts a value acceptable to time.ParseDuration)

