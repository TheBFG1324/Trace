package template_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"trace/package/utils/template"
)

func TestLoadJSON(t *testing.T) {
	tests := []struct {
		name           string
		jsonTemplate   map[string]interface{}
		taskParameters map[string]interface{}
		globalData     map[string]interface{}
		expectedJSON   string
		expectError    bool
	}{
		{
			name: "All placeholders resolved using task parameters",
			jsonTemplate: map[string]interface{}{
				"name":    "[[name]]",
				"age":     "[[age]]",
				"address": "[[address]]",
			},
			taskParameters: map[string]interface{}{
				"name":    "John",
				"age":     30,
				"address": "123 Main St",
			},
			globalData:   map[string]interface{}{},
			expectedJSON: `{"address":"123 Main St","age":30,"name":"John"}`,
			expectError:  false,
		},
		{
			name: "Placeholders resolved using both task parameters and global data",
			jsonTemplate: map[string]interface{}{
				"name":    "[[name]]",
				"age":     "[[age]]",
				"address": "[[address]]",
			},
			taskParameters: map[string]interface{}{
				"name": "John",
			},
			globalData: map[string]interface{}{
				"age":     30,
				"address": "123 Main St",
			},
			expectedJSON: `{"address":"123 Main St","age":30,"name":"John"}`,
			expectError:  false,
		},
		{
			name: "Unresolved placeholders return an error",
			jsonTemplate: map[string]interface{}{
				"name": "[[name]]",
				"age":  "[[age]]",
			},
			taskParameters: map[string]interface{}{
				"name": "John",
			},
			globalData:   map[string]interface{}{},
			expectedJSON: "",
			expectError:  true,
		},
		{
			name: "Nested placeholders resolved",
			jsonTemplate: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "[[name]]",
					"info": map[string]interface{}{
						"age": "[[age]]",
					},
				},
			},
			taskParameters: map[string]interface{}{
				"name": "John",
			},
			globalData: map[string]interface{}{
				"age": 30,
			},
			expectedJSON: `{"user":{"info":{"age":30},"name":"John"}}`,
			expectError:  false,
		},
		{
			name: "Array placeholders resolved",
			jsonTemplate: map[string]interface{}{
				"items": []interface{}{
					"[[item1]]",
					"[[item2]]",
				},
			},
			taskParameters: map[string]interface{}{
				"item1": "value1",
			},
			globalData: map[string]interface{}{
				"item2": "value2",
			},
			expectedJSON: `{"items":["value1","value2"]}`,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := template.LoadJSON(tt.jsonTemplate, tt.taskParameters, tt.globalData)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError {
				// Unmarshal to compare JSON structure
				var resultMap, expectedMap map[string]interface{}
				err = json.Unmarshal([]byte(result), &resultMap)
				if err != nil {
					t.Fatalf("Error unmarshaling result: %v", err)
				}

				err = json.Unmarshal([]byte(tt.expectedJSON), &expectedMap)
				if err != nil {
					t.Fatalf("Error unmarshaling expected JSON: %v", err)
				}

				if !reflect.DeepEqual(resultMap, expectedMap) {
					t.Errorf("Expected: %v, got: %v", expectedMap, resultMap)
				}
			}
		})
	}
}
