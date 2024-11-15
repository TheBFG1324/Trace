package main

import (
	"trace/package/scheduler"
	"trace/package/parser"
	"fmt"
)

func main() {
	input := `
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
`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	parentRequest := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	// Simulate the execution
	fmt.Println("Starting Execution:")
	scheduler.RunParentRequest(parentRequest)
}
