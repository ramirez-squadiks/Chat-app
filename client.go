package main

import (
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    client := &Client{conn: conn, send: make(chan []byte)}
    hub.register(client) // Регистрация клиента в хабе

    go client.writeMessages() // Запускаем горутину для отправки сообщений
    client.readMessages() // Чтение сообщений от клиента
}

func (c *Client) readMessages() {
    defer func() {
        hub.unregister(c) // Удаляем клиента из хаба
        c.conn.Close()
    }()
    for {
        _, msg, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        hub.broadcast <- msg // Отправляем сообщение в хаб
    }
}

func (c *Client) writeMessages() {
    defer c.conn.Close()
    for msg := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            break
        }
    }
}
