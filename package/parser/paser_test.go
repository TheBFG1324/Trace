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
    DATA globalDataAdd TYPE String ;
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
					"globalDataAdd": {
						DataName:     "globalDataAdd",
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
						AgentName: "Agent1",
						DataPermissions: map[string][]string{
							"data1": {"READ", "WRITE"},
						},
					},
					"Agent2": {
						AgentName: "Agent2",
						DataPermissions: map[string][]string{
							"data2":         {"READ"},
							"globalDataAdd": {"ADD"},
						},
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
						AgentName: "Worker",
						DataPermissions: map[string][]string{
							"config":  {"READ"},
							"results": {"WRITE"},
						},
					},
					"Reporter": {
						AgentName: "Reporter",
						DataPermissions: map[string][]string{
							"results": {"READ"},
						},
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
		{
			name: "Test 3: Nested RUNSEQ and RUNCON blocks",
			input: `
START
    DATA sharedData TYPE Map ;
    DATA tempData1 TYPE String ;
    DATA tempData2 TYPE String ;

    PERM AGENT MainAgent DATA sharedData ACCESS WRITE ;
    PERM AGENT Worker1 DATA sharedData ACCESS READ ;
    PERM AGENT Worker1 DATA tempData1 ACCESS WRITE ;
    PERM AGENT Worker2 DATA sharedData ACCESS READ ;
    PERM AGENT Worker2 DATA tempData2 ACCESS WRITE ;

    RUNSEQ {
        TASK Initialize AGENT MainAgent PARAMETERS (output=sharedData) ;
        RUNCON {
            RUNSEQ {
                TASK Process1 AGENT Worker1 PARAMETERS (input=sharedData, output=tempData1) ;
                TASK Finalize1 AGENT Worker1 PARAMETERS (input=tempData1) ;
            }
            RUNSEQ {
                TASK Process2 AGENT Worker2 PARAMETERS (input=sharedData, output=tempData2) ;
                TASK Finalize2 AGENT Worker2 PARAMETERS (input=tempData2) ;
            }
        }
        TASK Aggregate AGENT MainAgent PARAMETERS (input1=tempData1, input2=tempData2) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"sharedData": {
						DataName:     "sharedData",
						DataType:     "Map",
						InitialValue: "",
					},
					"tempData1": {
						DataName:     "tempData1",
						DataType:     "String",
						InitialValue: "",
					},
					"tempData2": {
						DataName:     "tempData2",
						DataType:     "String",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"MainAgent": {
						AgentName: "MainAgent",
						DataPermissions: map[string][]string{
							"sharedData": {"WRITE"},
						},
					},
					"Worker1": {
						AgentName: "Worker1",
						DataPermissions: map[string][]string{
							"sharedData": {"READ"},
							"tempData1":  {"WRITE"},
						},
					},
					"Worker2": {
						AgentName: "Worker2",
						DataPermissions: map[string][]string{
							"sharedData": {"READ"},
							"tempData2":  {"WRITE"},
						},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "Initialize",
								AgentName: "MainAgent",
								Parameters: map[string]string{
									"output": "sharedData",
								},
							},
							&RunConBlock{
								Statements: map[string]interface{}{
									"RUNSEQ_0": &RunSeqBlock{
										Statements: []interface{}{
											&Task{
												TaskName:  "Process1",
												AgentName: "Worker1",
												Parameters: map[string]string{
													"input":  "sharedData",
													"output": "tempData1",
												},
											},
											&Task{
												TaskName:  "Finalize1",
												AgentName: "Worker1",
												Parameters: map[string]string{
													"input": "tempData1",
												},
											},
										},
									},
									"RUNSEQ_1": &RunSeqBlock{
										Statements: []interface{}{
											&Task{
												TaskName:  "Process2",
												AgentName: "Worker2",
												Parameters: map[string]string{
													"input":  "sharedData",
													"output": "tempData2",
												},
											},
											&Task{
												TaskName:  "Finalize2",
												AgentName: "Worker2",
												Parameters: map[string]string{
													"input": "tempData2",
												},
											},
										},
									},
								},
							},
							&Task{
								TaskName:  "Aggregate",
								AgentName: "MainAgent",
								Parameters: map[string]string{
									"input1": "tempData1",
									"input2": "tempData2",
								},
							},
						},
					},
				},
			},
		},

		{
			name: "Test 4: Complex permissions and data types",
			input: `
START
    DATA userProfiles TYPE List ;
    DATA processedData TYPE Dict ;
    DATA finalReport TYPE String ;

    PERM AGENT DataCollector DATA userProfiles ACCESS ADD ;
    PERM AGENT DataProcessor DATA userProfiles ACCESS READ ;
    PERM AGENT DataProcessor DATA processedData ACCESS WRITE ;
    PERM AGENT Analyst DATA processedData ACCESS READ ;
    PERM AGENT Analyst DATA finalReport ACCESS WRITE ;

    RUNSEQ {
        TASK CollectData AGENT DataCollector PARAMETERS (output=userProfiles) ;
        TASK ProcessData AGENT DataProcessor PARAMETERS (input=userProfiles, output=processedData) ;
        TASK Analyze AGENT Analyst PARAMETERS (input=processedData, output=finalReport) ;
        TASK Publish AGENT Analyst PARAMETERS (report=finalReport) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"userProfiles": {
						DataName:     "userProfiles",
						DataType:     "List",
						InitialValue: "",
					},
					"processedData": {
						DataName:     "processedData",
						DataType:     "Dict",
						InitialValue: "",
					},
					"finalReport": {
						DataName:     "finalReport",
						DataType:     "String",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"DataCollector": {
						AgentName: "DataCollector",
						DataPermissions: map[string][]string{
							"userProfiles": {"ADD"},
						},
					},
					"DataProcessor": {
						AgentName: "DataProcessor",
						DataPermissions: map[string][]string{
							"userProfiles":  {"READ"},
							"processedData": {"WRITE"},
						},
					},
					"Analyst": {
						AgentName: "Analyst",
						DataPermissions: map[string][]string{
							"processedData": {"READ"},
							"finalReport":   {"WRITE"},
						},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "CollectData",
								AgentName: "DataCollector",
								Parameters: map[string]string{
									"output": "userProfiles",
								},
							},
							&Task{
								TaskName:  "ProcessData",
								AgentName: "DataProcessor",
								Parameters: map[string]string{
									"input":  "userProfiles",
									"output": "processedData",
								},
							},
							&Task{
								TaskName:  "Analyze",
								AgentName: "Analyst",
								Parameters: map[string]string{
									"input":  "processedData",
									"output": "finalReport",
								},
							},
							&Task{
								TaskName:  "Publish",
								AgentName: "Analyst",
								Parameters: map[string]string{
									"report": "finalReport",
								},
							},
						},
					},
				},
			},
		},

		{
			name: "Test 5: Deeply nested RUNCON within RUNSEQ blocks",
			input: `
START
    DATA initialData TYPE String ;
    DATA intermediateData TYPE String ;
    DATA finalData TYPE String ;

    PERM AGENT Starter DATA initialData ACCESS WRITE ;
    PERM AGENT MiddleMan DATA initialData ACCESS READ ;
    PERM AGENT MiddleMan DATA intermediateData ACCESS WRITE ;
    PERM AGENT Finisher DATA intermediateData ACCESS READ ;
    PERM AGENT Finisher DATA finalData ACCESS WRITE ;

    RUNSEQ {
        TASK StartProcess AGENT Starter PARAMETERS (output=initialData) ;
        RUNSEQ {
            RUNCON {
                TASK IntermediateStep1 AGENT MiddleMan PARAMETERS (input=initialData, output=intermediateData) ;
                TASK IntermediateStep2 AGENT MiddleMan PARAMETERS (input=initialData, output=intermediateData) ;
            }
            TASK MergeData AGENT MiddleMan PARAMETERS (input=intermediateData) ;
        }
        TASK Finalize AGENT Finisher PARAMETERS (input=intermediateData, output=finalData) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"initialData": {
						DataName:     "initialData",
						DataType:     "String",
						InitialValue: "",
					},
					"intermediateData": {
						DataName:     "intermediateData",
						DataType:     "String",
						InitialValue: "",
					},
					"finalData": {
						DataName:     "finalData",
						DataType:     "String",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"Starter": {
						AgentName: "Starter",
						DataPermissions: map[string][]string{
							"initialData": {"WRITE"},
						},
					},
					"MiddleMan": {
						AgentName: "MiddleMan",
						DataPermissions: map[string][]string{
							"initialData":      {"READ"},
							"intermediateData": {"WRITE"},
						},
					},
					"Finisher": {
						AgentName: "Finisher",
						DataPermissions: map[string][]string{
							"intermediateData": {"READ"},
							"finalData":        {"WRITE"},
						},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "StartProcess",
								AgentName: "Starter",
								Parameters: map[string]string{
									"output": "initialData",
								},
							},
							&RunSeqBlock{
								Statements: []interface{}{
									&RunConBlock{
										Statements: map[string]interface{}{
											"IntermediateStep1": &Task{
												TaskName:  "IntermediateStep1",
												AgentName: "MiddleMan",
												Parameters: map[string]string{
													"input":  "initialData",
													"output": "intermediateData",
												},
											},
											"IntermediateStep2": &Task{
												TaskName:  "IntermediateStep2",
												AgentName: "MiddleMan",
												Parameters: map[string]string{
													"input":  "initialData",
													"output": "intermediateData",
												},
											},
										},
									},
									&Task{
										TaskName:  "MergeData",
										AgentName: "MiddleMan",
										Parameters: map[string]string{
											"input": "intermediateData",
										},
									},
								},
							},
							&Task{
								TaskName:  "Finalize",
								AgentName: "Finisher",
								Parameters: map[string]string{
									"input":  "intermediateData",
									"output": "finalData",
								},
							},
						},
					},
				},
			},
		},

		{
			name: "Test 6: Multiple agents and data dependencies",
			input: `
START
    DATA rawData TYPE String ;
    DATA cleanedData TYPE String ;
    DATA analyzedData TYPE String ;
    DATA report TYPE String ;

    PERM AGENT Scraper DATA rawData ACCESS WRITE ;
    PERM AGENT Cleaner DATA rawData ACCESS READ ;
    PERM AGENT Cleaner DATA cleanedData ACCESS WRITE ;
    PERM AGENT Analyst DATA cleanedData ACCESS READ ;
    PERM AGENT Analyst DATA analyzedData ACCESS WRITE ;
    PERM AGENT Reporter DATA analyzedData ACCESS READ ;
    PERM AGENT Reporter DATA report ACCESS WRITE ;

    RUNSEQ {
        TASK ScrapeData AGENT Scraper PARAMETERS (output=rawData) ;
        TASK CleanData AGENT Cleaner PARAMETERS (input=rawData, output=cleanedData) ;
        TASK AnalyzeData AGENT Analyst PARAMETERS (input=cleanedData, output=analyzedData) ;
        TASK GenerateReport AGENT Reporter PARAMETERS (input=analyzedData, output=report) ;
    }
END
`,
			expected: &ParentRequest{
				GlobalData: map[string]*Data{
					"rawData": {
						DataName:     "rawData",
						DataType:     "String",
						InitialValue: "",
					},
					"cleanedData": {
						DataName:     "cleanedData",
						DataType:     "String",
						InitialValue: "",
					},
					"analyzedData": {
						DataName:     "analyzedData",
						DataType:     "String",
						InitialValue: "",
					},
					"report": {
						DataName:     "report",
						DataType:     "String",
						InitialValue: "",
					},
				},
				Permissions: map[string]*Permission{
					"Scraper": {
						AgentName: "Scraper",
						DataPermissions: map[string][]string{
							"rawData": {"WRITE"},
						},
					},
					"Cleaner": {
						AgentName: "Cleaner",
						DataPermissions: map[string][]string{
							"rawData":     {"READ"},
							"cleanedData": {"WRITE"},
						},
					},
					"Analyst": {
						AgentName: "Analyst",
						DataPermissions: map[string][]string{
							"cleanedData":  {"READ"},
							"analyzedData": {"WRITE"},
						},
					},
					"Reporter": {
						AgentName: "Reporter",
						DataPermissions: map[string][]string{
							"analyzedData": {"READ"},
							"report":       {"WRITE"},
						},
					},
				},
				Statements: []interface{}{
					&RunSeqBlock{
						Statements: []interface{}{
							&Task{
								TaskName:  "ScrapeData",
								AgentName: "Scraper",
								Parameters: map[string]string{
									"output": "rawData",
								},
							},
							&Task{
								TaskName:  "CleanData",
								AgentName: "Cleaner",
								Parameters: map[string]string{
									"input":  "rawData",
									"output": "cleanedData",
								},
							},
							&Task{
								TaskName:  "AnalyzeData",
								AgentName: "Analyst",
								Parameters: map[string]string{
									"input":  "cleanedData",
									"output": "analyzedData",
								},
							},
							&Task{
								TaskName:  "GenerateReport",
								AgentName: "Reporter",
								Parameters: map[string]string{
									"input":  "analyzedData",
									"output": "report",
								},
							},
						},
					},
				},
			},
		},
		// ... [Other test cases remain the same, adjust inputs and expected outputs if necessary]
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
