package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Bind(req *http.Request, out interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, out)
}
