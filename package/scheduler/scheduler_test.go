package scheduler

import (
	"testing"
	"trace/package/parser"
)

// TestRunParentRequest tests the RunParentRequest function.
func TestRunParentRequest(t *testing.T) {
	input := `
START
    DATA data1 TYPE String VALUE "Initial Data" ;
    DATA data2 TYPE String ;
    DATA globalDataAdd TYPE String ;
    DATA globalData TYPE string ;

    PERM AGENT Agent1 DATA data1 ACCESS READ, WRITE ;
    PERM AGENT Agent2 DATA data2 ACCESS READ ;
    PERM AGENT Agent2 DATA globalDataAdd ACCESS ADD, WRITE ;

    RUNSEQ {
        TASK FetchData AGENT Agent1 PARAMETERS (source="DB", output=data1) ;
        RUNCON {
            TASK ProcessData AGENT Agent2 PARAMETERS (input=data1, output=data2) ;
            TASK LogData AGENT Agent3 PARAMETERS (input=data1) ;
        }
        TASK SaveData AGENT Agent4 PARAMETERS (input=data2) ;
    }
END
`
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	parentRequest := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parser errors:\n%v", p.Errors())
	}

	success := RunParentRequest(parentRequest)
	if !success {
		t.Fatalf("RunParentRequest returned false")
	}
}
