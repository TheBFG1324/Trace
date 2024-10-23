
package template

import (
	"encoding/json"
	"strings"
)

// LoadJSON builds the final JSON string with proper data types
func LoadJSON(jsonTemplate map[string]interface{}, taskParameters, globalData map[string]interface{}) (string, error) {
	templateCopy := deepCopyMap(jsonTemplate)

	replacePlaceholders(templateCopy, taskParameters, globalData)

	finalJson, err := json.Marshal(templateCopy)
	if err != nil {
		return "", err
	}

	return string(finalJson), nil
}

// Helper function to recursively replace placeholders
func replacePlaceholders(data map[string]interface{}, taskParameters, globalData map[string]interface{}) {
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
					data[key] = nil 
				}
			}
		case map[string]interface{}:
			replacePlaceholders(v, taskParameters, globalData)
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					replacePlaceholders(itemMap, taskParameters, globalData)
				} else if itemStr, ok := item.(string); ok {
					if strings.HasPrefix(itemStr, "[[") && strings.HasSuffix(itemStr, "]]") {
						varName := strings.TrimSuffix(strings.TrimPrefix(itemStr, "[["), "]]")
						if val, ok := taskParameters[varName]; ok {
							v[i] = val
						} else if val, ok := globalData[varName]; ok {
							v[i] = val
						} else {
							v[i] = nil
						}
					}
				}
			}
		}
	}
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
