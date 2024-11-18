package template

import (
	"encoding/json"
	"errors"
	"strings"
)

// LoadJSON builds the final JSON string with proper data types. It ensures that all placeholders are filled; otherwise, it returns an error.
func LoadJSON(jsonTemplate map[string]interface{}, taskParameters, globalData map[string]interface{}) (string, error) {
	templateCopy := deepCopyMap(jsonTemplate)

	// Replace placeholders and check for unfilled placeholders
	if !replacePlaceholders(templateCopy, taskParameters, globalData) {
		return "", errors.New("unfilled placeholders found in the template")
	}

	finalJson, err := json.Marshal(templateCopy)
	if err != nil {
		return "", err
	}

	return string(finalJson), nil
}

// Helper function to recursively replace placeholders. Returns false if any placeholders remain unfilled.
func replacePlaceholders(data map[string]interface{}, taskParameters, globalData map[string]interface{}) bool {
	allPlaceholdersFilled := true

	for key, value := range data {
		switch v := value.(type) {
		case string:
			if strings.HasPrefix(v, "[[") && strings.HasSuffix(v, "]]") {
				varName := strings.TrimSuffix(strings.TrimPrefix(v, "[["), "]]")
				if val, ok := taskParameters[varName]; ok {
					data[key] = val
				} else if val, ok := globalData[varName]; ok {
					data[key] = val
				} else {
					// Placeholder remains unfilled
					allPlaceholdersFilled = false
				}
			}
		case map[string]interface{}:
			if !replacePlaceholders(v, taskParameters, globalData) {
				allPlaceholdersFilled = false
			}
		case []interface{}:
			for i, item := range v {
				switch itemVal := item.(type) {
				case map[string]interface{}:
					if !replacePlaceholders(itemVal, taskParameters, globalData) {
						allPlaceholdersFilled = false
					}
				case string:
					if strings.HasPrefix(itemVal, "[[") && strings.HasSuffix(itemVal, "]]") {
						varName := strings.TrimSuffix(strings.TrimPrefix(itemVal, "[["), "]]")
						if val, ok := taskParameters[varName]; ok {
							v[i] = val
						} else if val, ok := globalData[varName]; ok {
							v[i] = val
						} else {
							allPlaceholdersFilled = false
						}
					}
				}
			}
		}
	}

	return allPlaceholdersFilled
}

// Helper function to deep copy a map
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[key] = deepCopyMap(v)
		case []interface{}:
			copy[key] = deepCopySlice(v)
		default:
			copy[key] = v
		}
	}
	return copy
}

// Helper function to deep copy a slice
func deepCopySlice(original []interface{}) []interface{} {
	copy := make([]interface{}, len(original))
	for i, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			copy[i] = deepCopyMap(v)
		case []interface{}:
			copy[i] = deepCopySlice(v)
		default:
			copy[i] = v
		}
	}
	return copy
}
