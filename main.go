package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Retrieve API key from environment variables
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        log.Fatal("API key is not set in environment variables")
    }

    // // Read the topic from topic.txt
    // topicFile, err := os.Open("topic.txt")
    // if err != nil {
    //     log.Fatalf("Error opening topic file: %v", err)
    // }
    // defer topicFile.Close()

    // scanner := bufio.NewScanner(topicFile)
    // topicContent := ""
    // for scanner.Scan() {
    //     topicContent += scanner.Text() + "\n"
    // }
    // description := "Essay Topic: "
    // topicContent = description + topicContent

    // Read the essay from essay.txt
    essayFile, err := os.Open("essay.txt")
    if err != nil {
        log.Fatalf("Error opening essay file: %v", err)
    }
    defer essayFile.Close()

    scanner := bufio.NewScanner(essayFile)
    essayContent := ""
    for scanner.Scan() {
        essayContent += scanner.Text() + "\n"
    }

    // // Check for errors during scanning
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading from file: %v", err)
    }

    // API 1
    thread, err := CreateThread()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the ID of the thread
    fmt.Printf("Thread ID: %s\n", thread.ID)

    // API 2
    _, err = CreateThreadMessage(thread.ID, essayContent)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // fmt.Printf("Thread message created: %+v\n", threadMessage)

    // API 3
    _, err = CreateThreadRun(thread.ID, "asst_rfyem8YCJcbmaiexwlHLO4YX")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    // add delay for 5 seconds
    time.Sleep(5 * time.Second)
    
    // API 4
    // threadID := "thread_qRFXumnhkErUDecZjuFBbRHC"
    messages := GetMessages(thread.ID)

    print(messages)

    // API 5
    runs, err := ListThreadRuns(thread.ID)
    if err != nil {
        fmt.Println("Error:", err)
        print("Error:", err)
    }
    for _, run := range runs.Data {
        fmt.Printf("Run ID: %s, Status: %s\n", run.ID, run.Status)
    }

    

    
}
