// Package json converts the webapi records into json.
package json

import (
	"encoding/json"
	"fmt"

	"github.com/diegomagalhaes-dev/go-service/app/tooling/docs/webapi"
)

// Transform converts the collection of webapi records into json.
func Transform(records []webapi.Record) error {
	data, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return err
	}

	fmt.Print(string(data))

	return nil
}
