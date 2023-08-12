package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func JsonStdout(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(os.Stdout, string(b)); err != nil {
		return err
	}

	return nil
}
