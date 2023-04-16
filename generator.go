package gowin32gen

import (
	"encoding/json"
	"io"
)

func Generate(source io.Reader, dest io.Writer) error {
	decoder := json.NewDecoder(source)
	data := ApiFile{}
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	//_, err = dest.Write([]byte(fmt.Sprintf("%#v", data)))
	for _, c := range data.Constants {
		dest.Write([]byte(c.String()))
		dest.Write([]byte{'\n'})
	}

	for _, t := range data.Types {
		dest.Write([]byte(t.Generate()))
		dest.Write([]byte{'\n'})
	}

	return err
}
