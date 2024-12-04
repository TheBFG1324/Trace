package scheduler_test

import (
	"testing"
	"trace/package/logger"
	"trace/package/parser"
	"trace/package/scheduler"
)

// TestRunParentRequest tests the RunParentRequest function.
func TestRunParentRequest(t *testing.T) {
	input := `
START
    DATA origin TYPE String VALUE "Kansas" ;
    DATA destination TYPE String VALUE "California" ;
    DATA date TYPE String VALUE "2023-12-25" ;
    DATA flightInfo TYPE String ;

    DATA location TYPE String VALUE "Los Angeles" ;
    DATA hotelDate TYPE String VALUE "2023-12-25" ;
    DATA guests TYPE Int VALUE 2 ;
    DATA hotelInfo TYPE String ;

    DATA pickup TYPE String VALUE "Airport" ;
    DATA dropoff TYPE String VALUE "Hotel" ;
    DATA time TYPE String VALUE "14:00" ;
    DATA rideInfo TYPE String ;

    PERM AGENT FlightGetter DATA origin ACCESS READ ;
    PERM AGENT FlightGetter DATA destination ACCESS READ ;
    PERM AGENT FlightGetter DATA date ACCESS READ ;
    PERM AGENT FlightGetter DATA flightInfo ACCESS WRITE ;

    PERM AGENT RoomBooker DATA location ACCESS READ ;
    PERM AGENT RoomBooker DATA hotelDate ACCESS READ ;
    PERM AGENT RoomBooker DATA guests ACCESS READ ;
    PERM AGENT RoomBooker DATA hotelInfo ACCESS WRITE ;

    PERM AGENT UberScheduler DATA pickup ACCESS READ ;
    PERM AGENT UberScheduler DATA dropoff ACCESS READ ;
    PERM AGENT UberScheduler DATA time ACCESS READ ;
    PERM AGENT UberScheduler DATA rideInfo ACCESS WRITE ;

    RUNSEQ {
        TASK ScheduleFlight AGENT FlightGetter PARAMETERS (origin=origin, destination=destination, date=date, OUTPUT=flightInfo) ;
        TASK BookHotel AGENT RoomBooker PARAMETERS (location=location, date=hotelDate, guests=guests, OUTPUT=hotelInfo) ;
        TASK ScheduleRide AGENT UberScheduler PARAMETERS (pickup=pickup, dropoff=dropoff, time=time, OUTPUT=rideInfo) ;
    }
END
`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)
	parentRequest := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parser errors:\n%v", p.Errors())
	}

    l := logger.NewLogger()
	success := scheduler.RunParentRequest(parentRequest, l)
    l.PrintAllLogs()

	if !success {
		t.Fatalf("RunParentRequest returned false")
	}
}
