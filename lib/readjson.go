package lib

import (
	"encoding/json"
	"os"
)

func ReadJSONFromFile(filePath string) ([]ItemJsonStruct, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []ItemJsonStruct

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}
