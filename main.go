package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"

    "github.com/go-resty/resty/v2"
    "github.com/joho/godotenv"
)

const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
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

    // Open the essay.txt file
    file, err := os.Open("essay.txt") // Replace with the actual path to your .txt file
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer file.Close()

    // Read the file content
    scanner := bufio.NewScanner(file)
    essayContent := ""
    for scanner.Scan() {
        essayContent += scanner.Text() + "\n"
    }

    // Check for errors during scanning
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading from file: %v", err)
    }

    // Set up the resty client and make the API call
    client := resty.New()

    response, err := client.R().
        SetAuthToken(apiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model": "gpt-3.5-turbo",
            "messages": []interface{}{
                map[string]interface{}{
                    "role":    "system",
                    "content": "Please grade the following essay based on content, organization, grammar and sentence structure, and vocabulary and spelling. Each category should be scored from 1 to 5, with a total holistic score out of 20. Deduct 1 point in total score for essays significantly under the word count or not divided into paragraphs. Provide detailed scores and feedback for each category.",
                },
                map[string]interface{}{
                    "role":    "user",
                    "content": essayContent,
                },
            },
            "max_tokens": 1024, // Adjust this value as needed for the length of the response
        }).
        Post(apiEndpoint)

    // Check for errors in the API response
    if err != nil {
        log.Fatalf("Error while sending the request: %v", err)
    }

    // Decode the JSON response
    body := response.Body()

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        log.Fatalf("Error while decoding JSON response: %v", err)
    }

    // Extract the content from the JSON response
    choices := data["choices"].([]interface{})
    if len(choices) > 0 {
        content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
        fmt.Println("Graded Essay Content:", content)
    } else {
        fmt.Println("No content returned in the response.")
    }
}
