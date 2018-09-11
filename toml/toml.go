package toml

import (
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"strings"
)

// Decoder can take configurations encoded as TOML strings and decode them to
// config structs
type Decoder struct{}

func (d Decoder) Decode(r io.Reader, config interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	toml.Decode(string(data), config)
	return nil
}

// CanDecode returns true if this is a TOML file
func (d Decoder) CanDecode(path string) bool {
	return strings.HasSuffix(path, ".toml")
}
