# drone-extension-router

Drone extension server that is capable of routing different kinds of [extension requests](https://docs.drone.io/extensions/overview/) from the Drone server.

## Usage

In order to run the extension router, you must provide a `DRONE_SECRET` environment variable that is used to authenticate http requests to the extension.

```
$ DRONE_SECRET=mysecret ./drone-extension-router
```

The extension router also takes a configuration file path (`DRONE_CONFIG_FILE`) that is used to enable/configure various extensions.

```
$ DRONE_CONFIG_FILE=config.yaml DRONE_SECRET=mysecret ./drone-extension-router
```

The configuration format is as follows:
```
---
convert:
  pathschanged:
    enable: true
```

## Built-in Plugins

### Convert

|Plugin|Description|
|-|-|
|[pathschanged](https://github.com/meltwater/drone-convert-pathschanged)|Include/exclude pipelines and pipeline steps based on paths changed.|

## Testing

Run unit tests
```
make test
```

Run local docker container
```
make docker-run
```
