
<p align="center">
<img alt="Parrot logo" src="assets/parrot_banner.svg" height="150"/>
</p>

<p align="center">
  Parrot the <a href="http://poeditor.com/">POEditor</a> pull-through cache
</p>

<p align="center">
  <a href="https://github.com/UNIwise/parrot/releases/latest"><img src="https://img.shields.io/github/v/release/UNIwise/parrot"></a>
</p>
<p align="center">
  <img src="https://api.stage.eu.wiseflow.io/badges/v1/namespace/parrot/deployment/parrot?text=stage">
  <img src="https://api.eu.wiseflow.io/badges/v1/namespace/parrot/deployment/parrot?text=prod">
</p>

<hr/>

This cache makes it possible to pull app translations directly from [POEditor](http://poeditor.com/) instead of compiling them into your app. This means no more re-builds every time translations are updated!

# Features

-  **Battle tested**: The software is in active use on the WISEflow platform with high request rates daily.
-  **All the formats**: Parrot can provide all formats supported by POEditor.
-   **Cache choices**: Parrot comes with a Filesystem and Redis cache to facilitate single app deployments and highly distributed deployments.
-   **OpenAPI**: The Parrot API has been documented in OpenAPI specification which can be found in the [doc/](/docs) directory.
-  **Easy deployment**: A docker image and helm chart is provided.

# Configuration

The server can be configured with a yaml configuration file specified with the `--config` option. Every configuration can be overridden with an environment variable, where dots are replaced by underscores, such that `API_TOKEN=xxx` will set the `api.token` value

| key                            | description                                                                  | type     | default                      |
| ------------------------------ | ---------------------------------------------------------------------------- | -------- | ---------------------------- |
| server.port                    | port for the main http server                                                | int      | `80`                         |
| server.gracePeriod             | grace period for the http server to shutdown                                 | duration | `10s`                        |
| log.level                      | log level                                                                    | string   | `info`                       |
| log.format                     | format of the log. Can be "text" or "json"                                   | string   | `json`                       |
| cache.type                     | type of cache to use for translations                                        | string   | `filesystem`                 |
| cache.ttl                      | time to live for cache items                                                 | duration | `1h`                         |
| cache.renewalThreshold         | threshold at which the server will preemptively fetch a new translation      | duration | `30m`                        |
| cache.filesystem.dir           | directory of the filesystem cache                                            | string   | default user cache directory |
| cache.redis.mode               | mode of the redis connection to back the redis cache. "single" or "sentinel" | string   | `single`                     |
| cache.redis.address            | address of the redis server, in case the single mode is used                 | string   |
| cache.redis.username           | username to authenticate against redis                                       | string   |
| cache.redis.password           | password for redis authentication                                            | string   |
| cache.redis.maxRetries         | max retries for redis client to connect to redis. Set to -1 for infinity     | int      | `-1`                         |
| cache.redis.db                 | redis db index                                                               | int      | `1`                          |
| cache.redis.sentinel.master    | master name for sentinel setup                                               | string   |
| cache.redis.sentinel.addresses | list of sentinel addresses                                                   | []string |
| cache.redis.sentinel.password  | password for authenticating against sentinel instances                       | string   |
| prometheus.enabled             | enable prometheus metrics                                                    | boolean  | `true`                       |
| prometheus.path                | expose prometheus metrics under path                                         | string   | `/metrics`                   |
| prometheus.port                | port to expose the prometheus metrics under                                  | int      | `9090`                       |
| api.token                      | secret token to authenticating against poeditor                              | string   |
| aws.region                    | AWS Region | string | `eu-west-1` |
| aws.bucket.name               | AWS bucket name | string | 


# API specification

The REST API of Parrot is documented in the OpenAPI format. The specification file can be found here [docs/api.yml](docs/api.yml) and a Swagger UI is available here [uniwise.github.io/parrot](https://uniwise.github.io/parrot).

# License

Parrot is available under the Apache 2 license.

This project uses open source components which have additional licensing terms.
