package agent

import (
	"sync"
)

// Agent represents the interface that all agents should implement.
type Agent interface {
	GetID() int
	GetName() string
}

// BaseAgent provides a base implementation of the Agent interface.
type BaseAgent struct {
	ID           int
	Name         string
	AgentType    string
	Endpoint     string
	Reputation   float32
	Capabilities []string
	mu           sync.Mutex
}

// NewBaseAgent creates a new BaseAgent instance.
func NewBaseAgent(id int, name, agentType, endpoint string, capabilities []string) *BaseAgent {
	return &BaseAgent{
		ID:           id,
		Name:         name,
		AgentType:    agentType,
		Endpoint:     endpoint,
		Capabilities: capabilities,
		Reputation:   0.0,
	}
}

// GetID returns the agent's ID.
func (a *BaseAgent) GetID() int {
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
