package basic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func UnmarshalJson(data []byte, config interface{}, errorOnUnmatchedKeys bool) (err error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	if errorOnUnmatchedKeys {
		dec.DisallowUnknownFields()
	}
	err = dec.Decode(&config)
	return
}

type UnmatchedTomlKeysError struct {
	Keys []toml.Key
}

func (e *UnmatchedTomlKeysError) Error() string {
	return fmt.Sprintf("There are keys in the config file that do not match any field in the given struct: %v", e.Keys)
}

func UnmarshalToml(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	metaData, err := toml.Decode(string(data), config)
	if err == nil && len(metaData.Undecoded()) > 0 && errorOnUnmatchedKeys {
		return &UnmatchedTomlKeysError{Keys: metaData.Undecoded()}
	}
	return nil
}

func UnmarshalYaml(data []byte, config interface{}, errorOnUnmatchedKeys bool) (err error) {
	if errorOnUnmatchedKeys {
		dec := yaml.NewDecoder(bytes.NewBuffer(data))
		dec.KnownFields(true)
		return dec.Decode(config)
	}
	return yaml.Unmarshal(data, config)
}
