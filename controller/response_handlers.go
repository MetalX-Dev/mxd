package controller

import "github.com/MetalX-Dev/mxcommon/protocol"

func (ctx Context) handleCommandExecutionResponse(commandExecutionResponse *protocol.CommandExecutionResponse) bool {
	// TODO: Handle command execution response
	return false
}

func (ctx Context) handleFileOperationResponse(fileOperationResponse *protocol.FileOperationResponse) bool {
	// TODO: Handle file operation response
	return false
}

func (ctx Context) handleAgentResponse(agentResponse *protocol.AgentResponse) bool {
	// switch agentResponse.Payload.Type {
	switch {
	// case "CommandExecutionResponse":
	case agentResponse.Payload.CommandExecutionResponse != nil:
		return ctx.handleCommandExecutionResponse(agentResponse.Payload.CommandExecutionResponse)
		// case "FileOperationResponse":
	case agentResponse.Payload.FileOperationResponse != nil:
		return ctx.handleFileOperationResponse(agentResponse.Payload.FileOperationResponse)
	}
	return false
}
