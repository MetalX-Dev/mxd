package controller

import (
	"io"
	"log"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GetHealthResponse struct {
	Status string `json:"status"`
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	result := GetHealthResponse{
		Status: "OK",
	}
	response, err := sonic.Marshal(&result)
	if err != nil {
		log.Printf("Error marshalling response: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

type HostInfo struct {
	ID          string `json:"id"`
	ConnectedAt int64  `json:"connected_at"`
}
type GetHostListResponse struct {
	Count int `json:"count"`
	Hosts []HostInfo
}

func getHostList(w http.ResponseWriter, r *http.Request) {
	hosts := make([]HostInfo, len(pool))
	index := 0
	for _, v := range pool {
		hosts[index] = HostInfo{ID: v.id, ConnectedAt: v.connectedAt.UnixMilli()}
		index++
	}
	result := GetHostListResponse{
		Count: len(pool),
		Hosts: hosts[:],
	}
	response, err := sonic.Marshal(&result)
	if err != nil {
		log.Printf("Error marshalling response: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

type PostCmdExecRequest struct {
	HostId  string `json:"host_id"`
	Command string `json:"command"`
}
type PostCmdExecResponse struct {
	TaskId uint64 `json:"task_id"`
}

func postCmdExec(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var req PostCmdExecRequest
	err = sonic.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Error unmarshalling request: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	handleCommandExecutionRequest(req.HostId, req.Command)
}

func getCmdExecOutput(w http.ResponseWriter, r *http.Request) {

}

func RegisterApiServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log.Printf("Registering API server")
	r.Get("/api/v1/health", getHealth)
	r.Get("/api/v1/list", getHostList)
	r.Post("/api/v1/cmd_exec", postCmdExec)
	r.Get("/api/v1/cmd_exec/{task_id}", getCmdExecOutput)

	return r
}
