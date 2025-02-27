package controller

import "github.com/MetalX-Dev/mxcommon/protocol"

const VERSION = 0x0001

func handleControllerRequestPayload(hostId string, payload protocol.ControllerRequestPayload) bool {
	if ctx, ok := pool[hostId]; ok {
		ctx.seq += 1
		ctx.handleRequest(&protocol.ControllerRequest{
			Version: VERSION,
			ID:      ctx.seq,
			Payload: payload,
		})
		return true
	}

	return false
}

func handleCommandExecutionRequest(hostId string, command string) bool {
	return handleControllerRequestPayload(hostId, protocol.ControllerRequestPayload{
		// Type: "CommandExecutionRequest",
		CommandExecutionRequest: &protocol.CommandExecutionRequest{
			Command: command,
		},
	})
}
