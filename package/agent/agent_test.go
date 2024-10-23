package agent_test

import (
	"testing"
	"trace/package/agent"
	"reflect"
)
//TestAgentCreation tests the NewBaseAgent function.
func TestAgentCreation(t *testing.T) {
	
	jsonTemplate := map[string]interface{}{
		"taskName": "[[taskName]]",
		"age":      "[[age]]",
		"config": map[string]interface{}{
			"timeout": "[[timeout]]",
		},
	}

	agent := agent.NewBaseAgent(
		"AGT1557",
		"Example Agent",
		"BiometricAuthenticator",
		"http://example.com/api",
		jsonTemplate,
		[]string{"BiometricAuthentication-Signature", "BiometricAuthentication-Face"},
	)

	id := agent.GetID()
	name := agent.GetName()
	agentType := agent.GetAgentType()
	endpoint := agent.GetEndpoint()
	reputation := agent.GetReputation()
	capabilities := agent.GetCapabilities()
	jsonBody := agent.GetJsonBody()

	if id != agent.ID {
		t.Errorf("Expected response %s, but got %s", agent.ID, id)
	}

	if name != agent.Name {
		t.Errorf("Expected response %s, but got %s", agent.Name, name)
	}

	if agentType != agent.AgentType {
		t.Errorf("Expected response %s, but got %s", agent.AgentType, agentType)
	}

	if endpoint != agent.Endpoint {
		t.Errorf("Expected response %s, but got %s", agent.Endpoint, endpoint)
	}

	if reputation != agent.Reputation {
		t.Errorf("Expected response %v, but got %v", agent.Reputation, reputation)
	}

	if !reflect.DeepEqual(capabilities, agent.Capabilities) {
		t.Errorf("Expected response %v, but got %v", agent.Capabilities, capabilities)
	}

	if !reflect.DeepEqual(jsonBody, agent.JsonBody) {
		t.Errorf("Expected response %s, but got %s", agent.JsonBody, jsonBody)
	}
}