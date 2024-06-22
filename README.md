# goulp

goulp is a small set of Go utility packages for configuration loading,
structured logging, cron job registration, and JSON/YAML convenience helpers.
The repository is organized as independent packages; import only the package
you need.

The current codebase targets Go modules and declares `go 1.22` in `go.mod`.

## Packages

| Package | Purpose | Main exported API |
| --- | --- | --- |
| `confread` | Read one config file name from ordered candidate directories and unmarshal the first file found. | `Read`, `TryAllUnmarshalers`, `Unmarshalers`, `ErrNotFound` |
| `crontask` | Register and run keyed cron jobs on top of `robfig/cron/v3`. | `NewExecutor`, `JobMeta`, `Punch`, `RegValidator`, `CtxMaker`, registration and lookup errors |
| `jsonex` | Re-export a `json-iterator/go` configuration compatible with the standard library. | `Marshal`, `MarshalToString`, `MarshalIndent`, `Unmarshal`, `UnmarshalFromString`, `Get`, `NewEncoder`, `NewDecoder`, `Valid`, `MustMarshalToString` |
| `wlog` | Wrap `logrus` entries with context-aware fingerprints and local dev logging helpers. | `NewWLog`, `Common`, `ByCtx`, `ByCtxAndCache`, `ByCtxAndRemoveCache`, `Log`, `WLog`, `LDev`, `LInit`, `LExit` |
| `yaml` | Load YAML files or bytes with recursive `!include` support. | `LoadYAMLFile`, `LoadYAML`, `FileReader`, `IncludeTag`, `ErrIncludeCycle` |

## Install

```sh
go get github.com/bagaking/goulp
```

Import the package you need directly:

```go
import "github.com/bagaking/goulp/jsonex"
```

## Usage

### `confread`: read config from candidate directories

`confread.Read` builds candidate paths from the provided directories, reads the
first existing file, then tries YAML, TOML, and JSON unmarshalling in order.

```go
package main

import (
	"errors"
	"log"

	"github.com/bagaking/goulp/confread"
)

type AppConfig struct {
	Name string `yaml:"name" toml:"name" json:"name"`
}

func main() {
	var cfg AppConfig
	err := confread.Read("app.yaml", &cfg, "./config", ".")
	if errors.Is(err, confread.ErrNotFound) {
		log.Fatal("config file was not found")
	}
	if err != nil {
		log.Fatal(err)
	}
}
```

### `crontask`: register keyed cron jobs

`crontask.NewExecutor` accepts a `cron.Parser` and optional `logrus.Entry`.
Registered jobs are keyed, validated by the cron parser, and can also be
triggered manually by key.

```go
package main

import (
	"context"
	"log"

	"github.com/bagaking/goulp/crontask"
	"github.com/robfig/cron/v3"
)

func main() {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	executor := crontask.NewExecutor(parser, nil)

	err := executor.Register(context.Background(), crontask.JobMeta{
		Key:     "heartbeat",
		Crontab: "@every 1m",
		Punch: func(ctx context.Context) error {
			log.Println("heartbeat")
			return nil
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	executor.Start()
	defer executor.StopAndWait()

	if err := executor.TriggerByJobKey(context.Background(), "heartbeat"); err != nil {
		log.Fatal(err)
	}
}
```

### `jsonex`: standard-compatible JSON helpers

`jsonex` exposes `jsoniter.ConfigCompatibleWithStandardLibrary` helpers as
package-level functions.

```go
package main

import (
	"fmt"
	"log"

	"github.com/bagaking/goulp/jsonex"
)

func main() {
	payload := map[string]string{"name": "goulp"}

	raw, err := jsonex.MarshalToString(payload)
	if err != nil {
		log.Fatal(err)
	}

	var decoded map[string]string
	if err := jsonex.UnmarshalFromString(raw, &decoded); err != nil {
		log.Fatal(err)
	}

	fmt.Println(decoded["name"])
}
```

### `wlog`: context-aware `logrus` entries

The default logger can create entries from a context and attach fingerprints to
the `method_` and `finger_print_` fields.

```go
package main

import (
	"context"

	"github.com/bagaking/goulp/wlog"
)

func main() {
	ctx := context.Background()

	logger, ctx := wlog.ByCtxAndCache(ctx, "worker", "sync")
	logger.Info("started")

	wlog.ByCtx(ctx, "step").Info("running")
	wlog.LDev.Log("local").Debug("developer log")
}
```

### `yaml`: load YAML with `!include`

`yaml.LoadYAMLFile` reads a file from disk. `yaml.LoadYAML` accepts raw bytes, a
base directory for relative includes, and a custom file reader, which makes it
usable in tests or embedded file systems.

```go
package main

import (
	"log"

	"github.com/bagaking/goulp/yaml"
)

type Config struct {
	Name    string `yaml:"name"`
	Profile struct {
		User string `yaml:"user"`
	} `yaml:"profile"`
}

func main() {
	var cfg Config
	if err := yaml.LoadYAMLFile("config.yaml", &cfg); err != nil {
		log.Fatal(err)
	}
}
```

Example YAML:

```yaml
name: app
profile: !include profile.yaml
```

## Current Boundaries

- This module is a utility library, not a command-line tool or service.
- Packages are independent; there is no top-level `goulp` package API.
- `confread` chooses the first readable file from the candidate directories and
  tries all configured unmarshallers; it does not infer the decoder from the
  file extension.
- `crontask` wraps in-process cron scheduling only. It does not provide
  persistence, distributed locking, or cross-process coordination.
- `yaml.LoadYAML` processes `!include` directives through `yaml.Node` parsing
  before decoding into the destination value; unknown fields in the final
  output struct are currently not strictly rejected by this helper.
- recursive include cycles fail with `ErrIncludeCycle`.
- `wlog` is built on `logrus`; it does not replace `logrus` configuration for
  applications that need custom formatters, hooks, or outputs.

## Local Validation

```sh
git diff --check
git diff --cached --check
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
