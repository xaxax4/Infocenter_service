package handlers

import (
	"Infocenter/internal/services"
	"fmt"
	"io"
	"net/http"
)

var middleman *services.MiddleMan

func StartMiddleman(m *services.MiddleMan) {
	middleman = m
}

func InfocenterGETHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[InfocenterGETHandler] Called")
	topic := r.PathValue("topic")
	fmt.Printf("[InfocenterGETHandler] Topic: %s\n", topic)

	client := middleman.AddClient(topic)
	clientChan := client.Chat
	defer middleman.RemoveClient(topic, client)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	for {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			if msg.ID == -1 {
				_, err := fmt.Fprintf(w, "event: timeout\ndata: %s\n\n", msg.Content)
				if err != nil {
					return
				}
				flusher.Flush()
				return
			}
			_, err := fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Content)
			if err != nil {
				return
			}
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func InfocenterPOSTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[InfocenterPOSTHandler] Called")
	topic := r.PathValue("topic")
	fmt.Printf("[InfocenterPOSTHandler] Topic: %s\n", topic)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	middleman.SendToClients(topic, string(body))
	w.WriteHeader(http.StatusNoContent)
}
