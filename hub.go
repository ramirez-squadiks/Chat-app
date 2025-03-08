package main

import (
    "sync"
    "github.com/gorilla/websocket"
)

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

type Hub struct {
    clients   map[*Client]bool
    broadcast chan []byte
    mu        sync.Mutex
}

var hub = Hub{
    clients:   make(map[*Client]bool),
    broadcast: make(chan []byte),
}

func (h *Hub) run() {
    for {
        msg := <-h.broadcast
        h.mu.Lock()
        for client := range h.clients {
            select {
            case client.send <- msg:
            default:
                close(client.send)
                delete(h.clients, client)
            }
        }
        h.mu.Unlock()
    }
}

func (h *Hub) register(client *Client) {
    h.mu.Lock()
    h.clients[client] = true
    h.mu.Unlock()
}

func (h *Hub) unregister(client *Client) {
    h.mu.Lock()
    if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
    }
    h.mu.Unlock()
}
