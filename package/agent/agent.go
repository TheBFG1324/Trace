package agent

import (
	"sync"
	"fmt"
)

// BaseAgent provides a base implementation of the Agent interface.
type BaseAgent struct {
	ID           string
	Name         string
	AgentType    string
	Endpoint     string
	JsonBody     map[string]interface{}
	Reputation   float32
	Capabilities []string
	mu           sync.Mutex
}

// NewBaseAgent creates a new BaseAgent instance.
func NewBaseAgent(id string, name string, agentType string, endpoint string, jsonBody map[string]interface{}, capabilities []string) *BaseAgent {
	return &BaseAgent{
		ID:           id,
		Name:         name,
		AgentType:    agentType,
		Endpoint:     endpoint,
		JsonBody:     jsonBody,
		Capabilities: capabilities,
		Reputation:   0.0,
	}
}

// GetID returns the agent's ID.
func (a *BaseAgent) GetID() string {
	return a.ID
}

// GetName returns the agent's name.
func (a *BaseAgent) GetName() string {
	return a.Name
}

// GetAgentType returns the agent's type.
func (a *BaseAgent) GetAgentType() string {
	return a.AgentType
}

// GetEndpoint returns the agent's endpoint.
func (a *BaseAgent) GetEndpoint() string {
	return a.Endpoint
}

// GetReputation returns the agent's reputation.
func (a *BaseAgent) GetReputation() float32 {
	return a.Reputation
}

// GetCapabilities returns the agent's capabilities.
func (a *BaseAgent) GetCapabilities() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Capabilities
}

// GetJson returns the agent's json template.
func (a *BaseAgent) GetJsonBody() map[string]interface{} {
	return a.JsonBody
}

// GetMockAgents returns an array of mock base agents.
func GetMockAgents() []*BaseAgent {
	mockAgents := []*BaseAgent{
		NewBaseAgent("AG123", "FlightGetter", "Travel", "https://api.flightgetter.com",
			map[string]interface{}{
				"action": "search",
				"params": map[string]interface{}{
					"origin":      "[[origin]]",
					"destination": "[[destination]]",
					"date":        "[[date]]",
				},
			}, []string{"Search Flights", "Get Deals"}),
		NewBaseAgent("AG124", "RoomBooker", "Hospitality", "https://api.roombooker.com",
			map[string]interface{}{
				"action": "reserve",
				"params": map[string]interface{}{
					"location": "[[location]]",
					"date":     "[[date]]",
					"guests":   "[[guests]]",
				},
			}, []string{"Search Rooms", "Make Reservations"}),
		NewBaseAgent("AG125", "UberScheduler", "Transportation", "https://api.uberscheduler.com",
			map[string]interface{}{
				"action": "schedule",
				"params": map[string]interface{}{
					"pickup":   "[[pickup]]",
					"dropoff":  "[[dropoff]]",
					"time":     "[[time]]",
				},
			}, []string{"Schedule Ride", "Get ETA"}),
		NewBaseAgent("AG126", "WeatherChecker", "Utility", "https://api.weatherchecker.com",
			map[string]interface{}{
				"action": "get_weather",
				"params": map[string]interface{}{
					"location": "[[location]]",
					"date":     "[[date]]",
				},
			}, []string{"Get Weather", "Hourly Forecast"}),
		NewBaseAgent("AG127", "PackageTracker", "Logistics", "https://api.packagetracker.com",
			map[string]interface{}{
				"action": "track",
				"params": map[string]interface{}{
					"tracking_number": "[[tracking_number]]",
				},
			}, []string{"Track Package", "Delivery ETA"}),
	}
	return mockAgents
}


// SimulateLoadAgent returns predefined agents based on the given identifier type and value.
func SimulateLoadAgent(identifierType string, identifierValue string) *BaseAgent {
	mockAgents := GetMockAgents()

	for _, agent := range mockAgents {
		switch identifierType {
		case "ID":
			if agent.ID == identifierValue {
				return agent
			}
		case "Name":
			if agent.Name == identifierValue {
				return agent
			}
		case "AgentType":
			if agent.AgentType == identifierValue {
				return agent
			}
		default:
			fmt.Printf("Unknown identifier type: %s\n", identifierType)
			return nil
		}
	}

	fmt.Printf("Agent not found for %s: %s\n", identifierType, identifierValue)
	return nil
}