# Trace

Trace is a coordination platform built around the Artificial Intelligence Coordination Language (AICL). It provides a way for multiple agents (software services or components) to cooperate seamlessly while maintaining resilient interactions and robust trace logs. These logs enable you to audit and track any cooperation or data flow between agents.

# Overview

* Parsing and Executing AICL Scripts: Trace reads AICL programs (scripts) and orchestrates agents to accomplish tasks defined by the script.
* Resilient Cooperation: Agents can run tasks concurrently or sequentially. Each agent can read or update shared global data—subject to permission rules defined in the script.
* Trace Logging: Trace captures detailed logs of every step, making it easy to track changes to global data, see which agents performed each task, and debug any issues.

# AICL Script Structure
An AICL script typically follows this pattern:

## Global Data Declaration
```shell
DATA variableName TYPE SomeType VALUE "SomeValue" ;
```
Declares shared, global data accessible (and optionally updatable) by agents.
TYPE can be String, Int, Map, or other custom data types you define.
VALUE sets an initial value.

## Permissions
```shell
PERM AGENT SomeAgent DATA variableName ACCESS READ WRITE ;
```
Defines which agent can read or write which global data variable.
Example: PERM AGENT FlightGetter DATA flightInfo ACCESS WRITE ;
Means FlightGetter can write to the global variable flightInfo.

## Tasks
```shell
TASK TaskName AGENT AgentName PARAMETERS (param1=someGlobalData, OUTPUT=someOtherGlobalData) ;
```
AgentName references an agent’s name defined externally (e.g., a mock agent or a real microservice).
PARAMETERS define which inputs/outputs the agent uses.
If you want the agent to update a global variable, use an OUTPUT parameter:
```shell
PARAMETERS (origin=origin, destination=destination, OUTPUT=flightInfo)
```
This tells the system that the agent’s result will be stored in flightInfo.
If you want your parameters to pull info from an existing global variable, simply match the parameter name to the global data name:
```shell
PARAMETERS (origin=origin)
```
This means origin is read from the global variable origin.

## Blocks

RUNSEQ: A block of tasks that run sequentially.
RUNCON: A block of tasks that run concurrently (in parallel).
Nested blocks allow you to combine sequential and concurrent execution flows.
Example AICL Structure
```shell
START
DATA origin TYPE String VALUE "Kansas" ;
DATA destination TYPE String VALUE "California" ;
DATA date TYPE String VALUE "2023-12-25" ;
DATA flightInfo TYPE String ;

PERM AGENT FlightGetter DATA origin ACCESS READ ;
PERM AGENT FlightGetter DATA destination ACCESS READ ;
PERM AGENT FlightGetter DATA date ACCESS READ ;
PERM AGENT FlightGetter DATA flightInfo ACCESS WRITE ;

RUNSEQ {
    TASK ScheduleFlight AGENT FlightGetter PARAMETERS (origin=origin, destination=destination, date=date, OUTPUT=flightInfo) ;

    RUNCON {
        TASK SomeParallelTask1 AGENT SomeAgent1 PARAMETERS (input=flightInfo) ;
        TASK SomeParallelTask2 AGENT SomeAgent2 PARAMETERS (input=flightInfo) ;
    }

    TASK FinalSummary AGENT SomeAgent3 PARAMETERS (info=flightInfo) ;
}

END
```

## Enrolling Agents
Trace knows how to interact with agents that are “enrolled” in the system. Each agent typically has a JSON template describing how it consumes or produces data. Within this template, placeholders should match the AICL global data variable names, but bracketed with [[...]]. For instance:

{
"action": "reserve",
"params": {
"date": "[[date]]",
"location": "[[location]]",
"guests": "[[guests]]"
}
}

When the script is executed, Trace replaces placeholders like [[date]] with the actual values from the global data.

## Trace Logs
During script execution, Trace records:

Which tasks ran (including start time, finish time, agent ID).
What data was read and written (a full audit trail).
Concurrent execution logs, showing parallel tasks and data merges.
These logs can be used to reconstruct an audit trail of how each piece of data was produced or modified.

## Getting Started
Write an AICL script: Declare your global data, set permissions, define tasks in either sequential or concurrent blocks.
Register your agents: Implement or mock the agents that correspond to the agent names in the AICL script. Provide JSON templates with placeholders like [[variableName]].
Run the script: Use Trace’s executor to parse your script and automatically orchestrate the agent calls.
Review Logs: Trace logs all interactions so you can verify correctness and debug if needed.
With Trace, you can easily script high-level cooperative flows across multiple agents—ensuring resilience, auditability, and organizational clarity for complex coordination tasks.
