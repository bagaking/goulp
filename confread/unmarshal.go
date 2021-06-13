package confread

import (
	"github.com/BurntSushi/toml"
	"github.com/bagaking/goulp/jsonex"
	"gopkg.in/yaml.v2"
)

type Unmarshaler func(in []byte, out interface{}) (err error)

var Unmarshalers = []Unmarshaler{
	yaml.Unmarshal,
	toml.Unmarshal,
	jsonex.Unmarshal,
}

func TryAllUnmarshalers(source []byte, out interface{}) (err error) {
	for _, fnU := range Unmarshalers {
		if err = fnU(source, out); err == nil {
			break
		}
	}
	return
}
