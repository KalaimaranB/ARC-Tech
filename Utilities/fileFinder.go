package Utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// LoadFilenamesConfig loads the config file into a slice of maps or single map depending on structure
func LoadFilenamesConfig(filename string) ([]map[string]interface{}, error) {
	// Read the config file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the config file content into a slice of maps (array format)
	var configs []map[string]interface{}
	err = json.Unmarshal(data, &configs)
	if err != nil {
		// If the unmarshalling into a slice failed, try unmarshalling into a single object
		var singleConfig map[string]interface{}
		err = json.Unmarshal(data, &singleConfig)
		if err != nil {
			return nil, fmt.Errorf("could not parse config: %v", err)
		}
		// Return a slice containing the single object if it's not an array
		configs = append(configs, singleConfig)
	}

	return configs, nil
}

// SearchFileNames takes a list of strings and searches for the corresponding file path in filenames.json
func SearchFileNames(keys ...string) (string, error) {
	// Load the filenames config
	configs, err := LoadFilenamesConfig("filenames.json")
	if err != nil {
		return "", fmt.Errorf("could not load filenames.json: %v", err)
	}

	// Search through the config for the given keys
	for _, config := range configs {
		// Traverse through the map for each key in the list
		result, err := getFileName(config, keys)
		if err == nil {
			return result, nil
		}
	}

	// Return an error if no match is found
	return "", fmt.Errorf("no matching file found for the given keys: %v", keys)
}

// searchInMap recursively searches through the map using the provided keys
func getFileName(m map[string]interface{}, keys []string) (string, error) {
	if len(keys) == 0 {
		return "", fmt.Errorf("no more keys to search")
	}

	// Get the first key in the list
	key := keys[0]

	// Check if the key exists in the map
	if value, exists := m[key]; exists {
		// If this is the last key, return the value as string
		if len(keys) == 1 {
			// Ensure the value is a string or return an error
			if strValue, ok := value.(string); ok {
				return strValue, nil
			} else {
				return "", fmt.Errorf("value for '%s' is not a string", key)
			}
		}

		// If there are more keys to process, check if the value is a map and recursively search it
		if nestedMap, ok := value.(map[string]interface{}); ok {
			return getFileName(nestedMap, keys[1:])
		} else {
			return "", fmt.Errorf("value for '%s' is not a map", key)
		}
	}

	// If key doesn't exist, return an error
	return "", fmt.Errorf("key '%s' not found", key)
}
