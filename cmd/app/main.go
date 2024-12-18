package main

import (
	"fmt"
	"trace/package/logger"
	"trace/package/parser"
	"trace/package/scheduler"
)

func main() {
	input := `
START
    DATA origin TYPE String VALUE "Chicago" ;
    DATA destination TYPE String VALUE "New York" ;
    DATA date TYPE String VALUE "2024-05-15" ;
    DATA flightInfo TYPE String ;

    DATA location TYPE String VALUE "Manhattan" ;
    DATA hotelDate TYPE String VALUE "2024-05-16" ;
    DATA guests TYPE Int VALUE 2 ;
    DATA hotelInfo TYPE String ;

    DATA pickup TYPE String VALUE "Hotel" ;
    DATA dropoff TYPE String VALUE "Airport" ;
    DATA time TYPE String VALUE "08:30" ;
    DATA rideInfo TYPE String ;

    DATA weatherLocation TYPE String VALUE "Manhattan" ;
    DATA weatherDate TYPE String VALUE "2024-05-16" ;
    DATA weatherInfo TYPE String ;

    DATA trackingNumber TYPE String VALUE "XYZ-123" ;
    DATA packageStatus TYPE String ;

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

    PERM AGENT WeatherChecker DATA weatherLocation ACCESS READ ;
    PERM AGENT WeatherChecker DATA weatherDate ACCESS READ ;
    PERM AGENT WeatherChecker DATA weatherInfo ACCESS WRITE ;

    PERM AGENT PackageTracker DATA trackingNumber ACCESS READ ;
    PERM AGENT PackageTracker DATA packageStatus ACCESS WRITE ;

    RUNSEQ {
        TASK ScheduleFlight AGENT FlightGetter PARAMETERS (origin=origin, destination=destination, date=date, OUTPUT=flightInfo) ;

        RUNCON {
            RUNSEQ {
                TASK CheckWeather AGENT WeatherChecker PARAMETERS (location=weatherLocation, date=weatherDate, OUTPUT=weatherInfo) ;
                TASK BookHotel AGENT RoomBooker PARAMETERS (location=location, date=hotelDate, guests=guests, OUTPUT=hotelInfo) ;
            }
            RUNSEQ {
                TASK TrackPackage AGENT PackageTracker PARAMETERS (tracking_number=trackingNumber, OUTPUT=packageStatus) ;
                TASK PackageFollowUp AGENT PackageTracker PARAMETERS (tracking_number=trackingNumber) ;
            }
        }

        TASK ScheduleRide AGENT UberScheduler PARAMETERS (pickup=pickup, dropoff=dropoff, time=time, OUTPUT=rideInfo) ;
    }
END
`

	// Lex & parse the script
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	parentRequest := p.ParseProgram()

	// Check for parse errors
	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	// Create a logger
	lg := logger.NewLogger()

	// Run the parent request (the script)
	fmt.Println("Starting Execution:")
	success := scheduler.RunParentRequest(parentRequest, lg)

	// Print logs
	lg.PrintAllLogs()

	if !success {
		fmt.Println("Execution failed.")
	} else {
		fmt.Println("Execution succeeded.")
	}
}
