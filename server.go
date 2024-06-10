package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // 連接的客戶端
var broadcast = make(chan Message)           // 廣播訊息

// Message 聊天訊息結構
type Message struct {
    Username string `json:"username"`
    Message  string `json:"message"`
}

var History []Message

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func main() {
    // WebSocket 路由
    http.HandleFunc("/ws", handleConnections)
    http.HandleFunc("/history", handleHistory)

    // 啟動 goroutine 處理訊息
    go handleMessages()

    // 啟動伺服器
    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    // 升級 HTTP 連接至 WebSocket
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    // 將客戶端加入 map
    clients[ws] = true

    // 持續接收訊息
    for {
        var msg Message
        // 讀取 JSON 格式的訊息並解碼到 Message struct
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, ws)
            break
        }
        // 發送訊息至廣播 channel
        broadcast <- msg
        History = append(History, msg)
    }
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	json.NewEncoder(w).Encode(History)
}

func handleMessages() {
    for {
        // 從廣播 channel 接收訊息
        msg := <-broadcast
        // 將訊息發送給每個連接的客戶端
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}