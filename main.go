package main

import (
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
    // Use your API KEY here
    err := godotenv.Load()
	if err != nil{
		return
	}
    apiKey := os.Getenv("GPT_API_KEY") 
    client := resty.New()
    
    // 輸入題目
    var topics string
    fmt.Printf("Essay Topics in English: ")
    fmt.Scan(&topics)
    // 輸入文章
    var essay string
    fmt.Printf("Essay Content in English: ")
    fmt.Scan(&essay)

    var history = "Essay Topics: "+topics+"\nEssay Content: "+essay+"\n"

    for {
        response, err := client.R().
            SetAuthToken(apiKey).
            SetHeader("Content-Type", "application/json").
            SetBody(map[string]interface{}{
                "model":      "gpt-3.5-turbo",
                "messages":   []interface{}{
                    map[string]interface{}{
                        "role":    "system",
                        "content": "Please grade the following essay based on content, organization, grammar and sentence structure, and vocabulary and spelling. Each category should be scored from 1 to 5, with a total holistic score out of 20. Deduct 1 point in total score for essays significantly under the word count or not divided into paragraphs. Provide detailed scores and feedback for each category.",
                    },
                    map[string]interface{}{
                        "role":    "user",
                        "content": history,
                    },
                },
                "max_tokens": 1024,
            }).
            Post(apiEndpoint)

        if err != nil {
            log.Fatalf("Error while sending send the request: %v", err)
        }

        body := response.Body()

        var data map[string]interface{}
        err = json.Unmarshal(body, &data)
        if err != nil {
            fmt.Println("Error while decoding JSON response:", err)
            return
        }

        // 輸出回覆
        content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
        fmt.Println(content)
        history += "Graging Assistance: "+content+"\n"

        // 輸入對話
        var question string
        fmt.Printf("You: ")
        fmt.Scan(&question)
        if question == "exit" {
            fmt.Println("Prompt input closed.")
            return
        }
        history += "You: "+question+"\n"
    }

}