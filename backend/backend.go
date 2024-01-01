package main

import (
	"encoding/json"
	"net/http"
	"fmt"

	// "github.com/go-zeromq/zmq4/security/null"
)

type ThreadIDResponse struct {
    ThreadID string `json:"threadID"` // 使用JSON标签确保字段名称格式正确
}

type UserInputRequest struct {
    UserInput string `json:"userInput"`
    ThreadID  string `json:"threadID"` // 添加ThreadID字段
}

func getThreadID(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头部
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // 允许前端源
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 允许的方法
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // 允许的头部
	// 检查是否为预检请求
    if r.Method == "OPTIONS" {
        // 对于预检请求，返回适当的头部即可
        w.WriteHeader(http.StatusOK)
        return
    }

	thread, err := CreateThread()
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Thread ID: %s\n", thread.ID)
	
	response := ThreadIDResponse{
		ThreadID: thread.ID,
	}

	// 将响应结构体序列化为JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error marshalling JSON: %v", err)
        return
    }

    // 设置响应头部为JSON并发送响应
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

// 示例API处理函数
func handleThreadMessagesRequest(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头部
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // 允许前端源
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 允许的方法
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // 允许的头部

    // 检查是否为预检请求
    if r.Method == "OPTIONS" {
        // 对于预检请求，返回适当的头部即可
        w.WriteHeader(http.StatusOK)
        return
    }
	
	var request UserInputRequest

	// 解析请求正文到结构体
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// 现在可以处理UserInput作为一个字符串
    userInput := request.UserInput
    threadID := request.ThreadID
    fmt.Println("Received user input:", userInput)
    fmt.Println("Received thread ID:", threadID)

	response, err := SendMessages(threadID, userInput)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Printf("\n\nResponse type: %T\n", response)
    // // 构建示例响应
    // response := ThreadMessagesResponse{
	// 	Object: "list",
	// 	Data: []ThreadMessage{
	// 		{
	// 			ID:        "msg_gyAl1qVHmaO2sq40xeuLOqI5",
	// 			Object:    "thread.message",
	// 			CreatedAt: 1703830684,
	// 			ThreadID:  "thread_PwipxABtw3tjEZAXb312khtt",
	// 			Role:      "assistant",
	// 			Content: []Content{
	// 				{
	// 					Type: "text",
	// 					Text: ContentText{
	// 						Value:       "\"In this modern era, dominated by the influence of social media, emojis have gradually become an integral component of our daily lives.\"",
	// 						Annotations: []interface{}{},
	// 					},
	// 				},
	// 			},
	// 			FileIDs:     []string{},
	// 			AssistantID: nil,
	// 			RunID:       nil,
	// 			Metadata:    struct{}{},
	// 		},
	// 		{
	// 			ID:        "msg_gyAl1qVHmaO2sq40xeuLOqI5",
	// 			Object:    "thread.message",
	// 			CreatedAt: 1703830684,
	// 			ThreadID:  "thread_PwipxABtw3tjEZAXb312khtt",
	// 			Role:      "user",
	// 			Content: []Content{
	// 				{
	// 					Type: "text",
	// 					Text: ContentText{
	// 						Value:       "how to improve this sentence:\"In this time and age, with the social media taking the world by storm, “emoji” gradually becomes an indispensable part of our everyday lives.\"?",
	// 						Annotations: []interface{}{},
	// 					},
	// 				},
	// 			},
	// 			FileIDs:     []string{},
	// 			AssistantID: nil,
	// 			RunID:       nil,
	// 			Metadata:    struct{}{},
	// 		},
	// 	},
	// 	FirstID: "msg_gyAl1qVHmaO2sq40xeuLOqI5",
	// 	LastID:  "msg_WZAIdfaC8uhOfZgzLeu1ykgw",
	// 	HasMore: false,
	// }

    // 设置响应头为JSON
    w.Header().Set("Content-Type", "application/json")
    // 发送JSON响应
	jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }else{
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func main() {
    http.HandleFunc("/api/get-thread-id", getThreadID)
    http.HandleFunc("/api/thread-messages", handleThreadMessagesRequest)
	http.ListenAndServe(":8080", nil)
}
