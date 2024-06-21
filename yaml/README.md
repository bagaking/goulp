# goulp/yaml

`github.com/bagaking/goulp/yaml` loads YAML into Go structs and resolves
recursive `!include` tags before decoding. It is a package inside the `goulp`
module, not a standalone command.

## Install

```sh
go get github.com/bagaking/goulp
```

Import the package directly:

```go
import "github.com/bagaking/goulp/yaml"
```

## Load A File

`LoadYAMLFile` reads a YAML file from disk. Relative `!include` paths are
resolved from the directory of the file that contains the include.

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

Example input:

```yaml
name: app
profile: !include profile.yaml
```

`profile.yaml`:

```yaml
user: gogo
```

## Use A Custom Reader

`LoadYAML` accepts raw YAML bytes, a base directory for relative includes, and a
`FileReader`. Use this path when tests or embedded files should control how
include files are read.

```go
files := map[string][]byte{
	"profile.yaml": []byte("user: gogo\n"),
}

readFile := func(filename string) ([]byte, error) {
	data, ok := files[filename]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

var cfg Config
err := yaml.LoadYAML([]byte("name: app\nprofile: !include profile.yaml\n"), ".", &cfg, readFile)
```

## Boundaries

- Unknown YAML fields are rejected because the decoder enables
  `KnownFields(true)`.
- Includes may be nested; nested relative paths are resolved from the included
  file's directory.
- The package expands YAML includes before decoding. It does not watch files,
  merge multiple top-level config files, or provide environment interpolation.

## Validation

From the repository root:

```sh
go test ./yaml
go test ./...
```
