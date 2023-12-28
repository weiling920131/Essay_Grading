package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// ThreadResponse represents the JSON structure of the response from OpenAI
type ThreadResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Metadata  struct{}
}

// ContentText represents the text structure within content
type ContentText struct {
    Value       string `json:"value"`
    Annotations []interface{} `json:"annotations"` // Assuming annotations is an array of objects
}

type Content struct {
    Type string `json:"type"`
    Text ContentText `json:"text"`
}

// ThreadRunResponse represents the JSON structure of the response for a thread run
type ThreadRunResponse struct {
    ID          string `json:"id"`
    Object      string `json:"object"`
    CreatedAt   int64  `json:"created_at"`
    ThreadID    string `json:"thread_id"`
    AssistantID string `json:"assistant_id"`
    // Include any other fields you need
}

type ThreadMessageResponse struct {
    ID        string `json:"id"`
    Object    string `json:"object"`
    CreatedAt int64  `json:"created_at"`
    ThreadID  string `json:"thread_id"`
    Role      string `json:"role"`
    Content   []Content `json:"content"`
    FileIDs   []string `json:"file_ids"`
    AssistantID *string `json:"assistant_id"` // Assuming it can be null
    RunID     *string `json:"run_id"` // Assuming it can be null
    Metadata  struct{} `json:"metadata"`
}

// ThreadMessagesResponse represents the JSON structure of the response for thread messages
type ThreadMessagesResponse struct {
    Object   string          `json:"object"`
    Data     []ThreadMessage `json:"data"`
    FirstID  string          `json:"first_id"`
    LastID   string          `json:"last_id"`
    HasMore  bool            `json:"has_more"`
}

type ThreadMessage struct {
    ID          string `json:"id"`
    Object      string `json:"object"`
    CreatedAt   int64  `json:"created_at"`
    ThreadID    string `json:"thread_id"`
    Role        string `json:"role"`
    Content     []Content `json:"content"`
    FileIDs     []string `json:"file_ids"`
    AssistantID *string `json:"assistant_id"` // Assuming it can be null
    RunID       *string `json:"run_id"` // Assuming it can be null
    Metadata    struct{} `json:"metadata"`
}


// CreateThread creates a new thread using the OpenAI API
func CreateThread() (*ThreadResponse, error) {
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key not set")
    }

    url := "https://api.openai.com/v1/threads"

    client := &http.Client{}

    requestBody, err := json.Marshal(map[string]string{
        // Add necessary fields if needed
    })
    if err != nil {
        return nil, fmt.Errorf("error marshalling request body: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

	// println(string(body))

    var threadResponse ThreadResponse
    err = json.Unmarshal(body, &threadResponse)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %w", err)
    }

    return &threadResponse, nil
}

// CreateThreadMessage posts a message to an existing thread
func CreateThreadMessage(threadID, content string) (*ThreadMessageResponse, error) {
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not set")
	}

	url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", threadID)

	client := &http.Client{}

	requestBody, err := json.Marshal(map[string]string{
		"role":    "user",
		"content": content,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

	// println(string(body))

	var threadMessageResponse ThreadMessageResponse
	err = json.Unmarshal(body, &threadMessageResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &threadMessageResponse, nil
}

// CreateThreadRun creates a run in a specified thread
func CreateThreadRun(threadID, assistantID string) (*ThreadRunResponse, error) {
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key not set")
    }

    url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", threadID)

    client := &http.Client{}

    requestBody, err := json.Marshal(map[string]string{
        "assistant_id": assistantID,
    })
    if err != nil {
        return nil, fmt.Errorf("error marshalling request body: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

    var runResponse ThreadRunResponse
    err = json.Unmarshal(body, &runResponse)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %w", err)
    }

    return &runResponse, nil
}


// ListThreadMessages retrieves messages from a thread
func ListThreadMessages(threadID string) (*ThreadMessagesResponse, error) {
    apiKey := os.Getenv("GPT_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key not set")
    }

    url := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", threadID)

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

	fmt.Println("ListThreadMessages:")
	print(string(body))
	println("")
	println("")
	println("")
	
    var threadMessagesResponse ThreadMessagesResponse
    err = json.Unmarshal(body, &threadMessagesResponse)
    if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
    }
	
    return &threadMessagesResponse, nil
}

func HasValidAssistantResponse(response *ThreadMessagesResponse) bool {
    if len(response.Data) > 0 {
        firstMessage := response.Data[0]
        if firstMessage.Role == "assistant" {
            for _, content := range firstMessage.Content {
                if content.Type == "text" && content.Text.Value != "" {
                    return true
                }
            }
        }
    }
    return false
}

func GetMessages(threadID string) (*ThreadMessagesResponse) {
	time.Sleep(15 * time.Second)
	for {
		messageResponse, _ := ListThreadMessages(threadID)
		runs, err := ListThreadRuns(threadID)
		
		firstMessage := messageResponse.Data[0]
        if firstMessage.Role == "assistant" {
			for _, content := range firstMessage.Content {
				if content.Type == "text" && content.Text.Value != "" {
					if err != nil {
						fmt.Println("Error:", err)
						print("Error:", err)
					}
					for _, run := range runs.Data {
						fmt.Printf("Run ID: %s, Status: %s\n", run.ID, run.Status)
						if run.Status == "in_progress" {
							err := CancelThreadRun(threadID, run.ID)
							if err != nil {
								fmt.Println("Error canceling thread run:", err)
							} else {
								fmt.Println("Run canceled successfully")
							}
						}
					}

                    return messageResponse
                }
            }
        }
		secondMessage := messageResponse.Data[1]
        if secondMessage.Role == "assistant" {
			for _, content := range secondMessage.Content {
				if content.Type == "text" && content.Text.Value != "" {
					if err != nil {
						fmt.Println("Error:", err)
						print("Error:", err)
					}
					for _, run := range runs.Data {
						fmt.Printf("Run ID: %s, Status: %s\n", run.ID, run.Status)
						if run.Status == "in_progress" {
							err := CancelThreadRun(threadID, run.ID)
							if err != nil {
								fmt.Println("Error canceling thread run:", err)
							} else {
								fmt.Println("Run canceled successfully")
							}
						}
					}

                    return messageResponse
                }
            }
        }

		time.Sleep(5 * time.Second)
	}

}