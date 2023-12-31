package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

// ThreadRun represents an individual thread run
type ThreadRun struct {
    ID            string `json:"id"`
    Object        string `json:"object"`
    CreatedAt     int64  `json:"created_at"`
    AssistantID   string `json:"assistant_id"`
    ThreadID      string `json:"thread_id"`
    Status        string `json:"status"`
    StartedAt     int64  `json:"started_at"`
    ExpiresAt     *int64 `json:"expires_at"`
    CancelledAt   *int64 `json:"cancelled_at"`
    FailedAt      *int64 `json:"failed_at"`
    CompletedAt   *int64 `json:"completed_at"`
    LastError     *string `json:"last_error"`
    Model         string `json:"model"`
    Instructions  string `json:"instructions"`
    Tools         []string `json:"tools"`
    FileIDs       []string `json:"file_ids"`
    Metadata      struct{} `json:"metadata"`
}

// ThreadRunsResponse represents the JSON structure of the response for thread runs
type ThreadRunsResponse struct {
    Object   string       `json:"object"`
    Data     []ThreadRun  `json:"data"`
    FirstID  string       `json:"first_id"`
    LastID   string       `json:"last_id"`
    HasMore  bool         `json:"has_more"`
}

// ListThreadRuns retrieves runs from a thread
func ListThreadRuns(threadID string) (*ThreadRunsResponse, error) {
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key not set")
    }

    url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", threadID)

    client := &http.Client{}

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Add("Authorization", "Bearer "+apiKey)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("OpenAI-Beta", "assistants=v1")  // Add this header

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error performing request: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    var runsResponse ThreadRunsResponse
    err = json.Unmarshal(body, &runsResponse)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %w", err)
    }

    return &runsResponse, nil
}

// CancelThreadRun cancels a run in a specified thread
func CancelThreadRun(threadID, runID string) error {
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        return fmt.Errorf("API key not set")
    }

    url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs/%s/cancel", threadID, runID)

    client := &http.Client{}

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil)) // Empty body for POST request
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Add("Authorization", "Bearer "+apiKey)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("OpenAI-Beta", "assistants=v1")  // Add this header

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error performing request: %w", err)
    }
    defer resp.Body.Close()

    // Optionally, you can read and handle the response here if needed

    return nil
}
