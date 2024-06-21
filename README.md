# goulp

goulp is a collection of Go utility packages.

## Packages

- `confread`: read configuration files from candidate directories and unmarshal them with supported decoders.
- `crontask`: small wrappers around `robfig/cron` for registering, running, stopping, and manually triggering cron jobs.
- `jsonex`: `json-iterator/go` helpers configured for standard library compatibility.
- `wlog`: context-aware logging helpers built on `logrus`.
- `yaml`: YAML loading helpers, including support for `!include` tags.

## Install

```sh
go get github.com/bagaking/goulp
```

Import the package you need:

```go
import "github.com/bagaking/goulp/jsonex"
```

## Local Validation

```sh
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
