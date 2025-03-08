package main

import (
    "fmt"
    "net/http"
)

func main() {
    go hub.run() 
    http.HandleFunc("/ws", handleConnection) 
    fmt.Println("Сервер запущен на :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Ошибка при запуске сервера:", err)
    }
}
