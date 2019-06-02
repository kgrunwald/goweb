package framework

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// LoadYaml reads a YAML file and decodes the ontents into `out`. The function expects the `file` to exist in the
// directory indicated by the `CONFIG_DIR` environment variable.
func LoadYaml(file string, out interface{}) error {
	path := fmt.Sprintf("%s/config/%s", os.Getenv("CONFIG_DIR"), file)
	data, _ := ioutil.ReadFile(path)

	return yaml.Unmarshal([]byte(data), out)
}
