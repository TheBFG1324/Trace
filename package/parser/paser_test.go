package main

import (
	"reflect"
	"testing"
)

// TestParser tests the parsing of various AICL scripts with larger inputs.
func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ParentRequest
	}{
		{
			name: "Test 1: Complex AICL script",
			input: `
START
    DATA data1 TYPE String VALUE "Initial Data" ;
    DATA data2 TYPE String ;
    DATA globalData TYPE string ;

    PERM AGENT Agent1 DATA data1 ACCESS READ, WRITE ;
    PERM AGENT Agent2 DATA data2 ACCESS READ ;
    PERM AGENT Agent2 DATA globalDataAdd ACCESS ADD ;

    RUNSEQ {
        TASK FetchData AGENT Agent1 PARAMETERS (source="DB", output=data1) ;
        RUNCON {
            TASK ProcessData AGENT Agent2 PARAMETERS (input=data1, output=data2) ;
            TASK LogData AGENT Agent3 PARAMETERS (input=data1) ;
        }
        WAIT ProcessData ;
        TASK SaveData AGENT Agent4 PARAMETERS (input=data2) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"data1": {
						DataName:     "data1",
						DataType:     "String",
						InitialValue: "Initial Data",
					},
					"data2": {
						DataName:     "data2",
						DataType:     "String",
						InitialValue: "",
					},
					"globalData": {
						DataName:     "globalData",
						DataType:     "string",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"Agent1": {
						AgentName:   "Agent1",
						DataNames:   []string{"data1"},
						Permissions: []string{"READ", "WRITE"},
					},
					"Agent2": {
						AgentName:   "Agent2",
						DataNames:   []string{"data2", "globalDataAdd"},
						Permissions: []string{"READ", "ADD"},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "FetchData",
								AgentName: "Agent1",
								Parameters: map[string]string{
									"source": "DB",
									"output": "data1",
								},
							},
							&RunConBlock{
								Statements: map[string]interface{}{
									"ProcessData": &Task{
										TaskName:  "ProcessData",
										AgentName: "Agent2",
										Parameters: map[string]string{
											"input":  "data1",
											"output": "data2",
										},
									},
									"LogData": &Task{
										TaskName:  "LogData",
										AgentName: "Agent3",
										Parameters: map[string]string{
											"input": "data1",
										},
									},
								},
							},
							&WaitStatement{
								TaskNames: []string{"ProcessData"},
							},
							&Task{
								TaskName:  "SaveData",
								AgentName: "Agent4",
								Parameters: map[string]string{
									"input": "data2",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Test 2: Script with comments and multiple permissions",
			input: `
START
    DATA config TYPE JSON VALUE "{\"setting\": true}" ;
    DATA results TYPE Array ;

    PERM AGENT Worker DATA config ACCESS READ ;
    PERM AGENT Worker DATA results ACCESS WRITE ;
    PERM AGENT Reporter DATA results ACCESS READ ;

    RUNSEQ {
        TASK Setup AGENT Worker PARAMETERS (config=config) ;
        RUNCON {
            TASK Compute1 AGENT Worker PARAMETERS (input=config, output=results) ;
            TASK Compute2 AGENT Worker PARAMETERS (input=config, output=results) ;
        }
        WAIT Compute1, Compute2 ;
        TASK Report AGENT Reporter PARAMETERS (data=results) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"config": {
						DataName:     "config",
						DataType:     "JSON",
						InitialValue: `{"setting": true}`,
					},
					"results": {
						DataName:     "results",
						DataType:     "Array",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"Worker": {
						AgentName:   "Worker",
						DataNames:   []string{"config", "results"},
						Permissions: []string{"READ", "WRITE"},
					},
					"Reporter": {
						AgentName:   "Reporter",
						DataNames:   []string{"results"},
						Permissions: []string{"READ"},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "Setup",
								AgentName: "Worker",
								Parameters: map[string]string{
									"config": "config",
								},
							},
							&RunConBlock{
								Statements: map[string]interface{}{
									"Compute1": &Task{
										TaskName:  "Compute1",
										AgentName: "Worker",
										Parameters: map[string]string{
											"input":  "config",
											"output": "results",
										},
									},
									"Compute2": &Task{
										TaskName:  "Compute2",
										AgentName: "Worker",
										Parameters: map[string]string{
											"input":  "config",
											"output": "results",
										},
									},
								},
							},
							&WaitStatement{
								TaskNames: []string{"Compute1", "Compute2"},
							},
							&Task{
								TaskName:  "Report",
								AgentName: "Reporter",
								Parameters: map[string]string{
									"data": "results",
								},
							},
						},
					},
				},
			},
		},
		// ... [Other test cases remain the same, adjust inputs if necessary]
	}

	for _, test := range tests {
		l := NewLexer(test.input)
		p := NewParser(l)
		actual := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("Parser errors in test '%s': %v", test.name, p.Errors())
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Test '%s' failed:\nExpected: %+v\nGot: %+v", test.name, test.expected, actual)
		}
	}
}
