# AI Coordination Language Executor

**AI Coordination Language (AICL) Executor** is a Go-based application designed to parse, execute, and manage interactions between AI agents using a specialized scripting language called **AICL (AI Coordination Language)**. AICL facilitates orchestrating tasks for multiple agents, supporting concurrency, synchronization, and controlled data access through a central supervisor.

## Table of Contents

- [Introduction to AICL](#introduction-to-aicl)
- [Application Overview](#application-overview)
- [Folder Structure](#folder-structure)
- [How It Works](#how-it-works)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Executing AICL Scripts](#executing-aicl-scripts)
- [Components Breakdown](#components-breakdown)
  - [Parser](#parser)
  - [Models](#models)
  - [Executor](#executor)
  - [Supervisor](#supervisor)
  - [Agents](#agents)
  - [Web Server](#web-server)
  - [Front End](#front-end)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Introduction to AICL

### What is AICL?

**AI Coordination Language (AICL)** is a lightweight, declarative scripting language designed to facilitate communication and task orchestration between multiple AI agents. Using AICL, a **Parent AI** can issue instructions to agents, defining:

- **Concurrent Execution** (`RUNCON`): Tasks to be run in parallel.
- **Sequential Execution** (`RUNSEQ`): Tasks to be run one after the other.
- **Task Definitions** (`TASK`): Individual units of work executed by specific agents.
- **Waiting for Tasks** (`WAIT`): Synchronization points where execution pauses until specified tasks complete.
- **Controlled Data Access** (`PERM`): Permissions for agents to read, write, or execute data.

- After the successful execution of the AICL script, all final global data and logs will be handed back to the **Parent AI**. 
- Contributing agents may be allowed to access global data final state if they were given permissions

### AICL Example

```plaintext
START
    DATA data1 TYPE String VALUE "Initial Data" ;
    DATA data2 TYPE String ;

    PERM AGENT Agent1 DATA data1 ACCESS READ, WRITE ;
    PERM AGENT Agent2 DATA data2 ACCESS READ ;

    RUNSEQ {
        TASK FetchData AGENT Agent1 PARAMETERS (cmd="update", source="DB", output=data1) ;
        RUNCON {
            TASK ProcessData AGENT Agent2 PARAMETERS (input=data1, output=data2) ;
            TASK LogData AGENT Agent3 PARAMETERS (input=data1) ;
        }
        WAIT ProcessData ;
        TASK SaveData AGENT Agent4 PARAMETERS (input=data2) ;
    }
END
```