package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MetalX-Dev/mxcommon/protocol"
	"github.com/gorilla/websocket"
)

const (
	IntMsgSendContent     = iota // Send string content via websocket
	IntMsgCloseConnection        // Close the websocket connection gracefully, not implemented
)

// internal purpose for delivering events to websocket context
type InternalMessage struct {
	messageType int
	content     string
}

// Context for agent session
type Context struct {
	internalChan *chan InternalMessage
	id           string
	connectedAt  time.Time
	seq          uint64
}

// pool of agent sessions
var pool = make(map[string]*Context)

// handle agent response, parse and dispatch to response handlers
func (ctx Context) handleResponse(message string) {
	log.Printf("Received response from host %s: %s", ctx.id, message)

	agentResponse, err := protocol.ParseAgentResponse(string(message))
	if err != nil {
		log.Printf("Error parsing message: %s", err)
	}

	ctx.handleAgentResponse(agentResponse)
}

// send ControllerRequest to agent
func (ctx Context) handleRequest(request *protocol.ControllerRequest) {
	body := request.String()
	log.Printf("Sending request to host %s: %s", ctx.id, body)

	*ctx.internalChan <- InternalMessage{
		messageType: IntMsgSendContent,
		content:     body,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handle WebSocket connection
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	hostId := r.URL.Query().Get("host_id")
	if hostId == "" {
		http.Error(w, "Missing host_id", http.StatusBadRequest)
		return
	}

	log.Printf("Handling WebSocket connection from %s", r.RemoteAddr)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	ctxConn := make(chan InternalMessage)
	ctx := &Context{
		internalChan: &ctxConn,
		id:           hostId,
		connectedAt:  time.Now(),
		seq:          0,
	}
	pool[hostId] = ctx

	defer func() {
		ws.Close()
		delete(pool, hostId)
		log.Printf("Closed WebSocket connection from %s", r.RemoteAddr)
	}()

	go func() {
		for message := range ctxConn {
			switch message.messageType {
			case IntMsgSendContent:
				err := ws.WriteMessage(websocket.TextMessage, []byte(message.content))
				if err != nil {
					log.Printf("Error writing message: %s", err)
				}
			default:
				log.Printf("Unexpected message type: %d", message.messageType)
			}
		}
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %s", err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			ctx.handleResponse(string(message))
		default:
			log.Printf("Received unknown message type: %d", messageType)
		}
	}
}

// Start the controller server
func StartServer(port int) {
	api_server := RegisterApiServer()
	api_server.Handle("/ws", http.HandlerFunc(handleWebSocket))
	http.ListenAndServe(fmt.Sprintf(":%d", port), api_server)
}
