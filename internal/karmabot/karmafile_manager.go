package karmabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func WriteKarma() error {
	karmafile := os.Getenv("KARMAFILE")
	b, err := json.Marshal(&Karma)
	if err == nil {
		err := ioutil.WriteFile(karmafile, b, 0644)
		if err != nil {
			return fmt.Errorf("Error: Couldn't save karma to %s", karmafile)
		}
	}
	return nil
}
