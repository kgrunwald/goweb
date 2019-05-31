package framework

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYaml(file string, out interface{}) error {
	path := fmt.Sprintf("%s/config/%s", os.Getenv("CONFIG_DIR"), file)
	data, _ := ioutil.ReadFile(path)

	return yaml.Unmarshal([]byte(data), out)
}
