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

    response, err := client.R().
        SetAuthToken(apiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":      "gpt-3.5-turbo",
            "messages":   []interface{}{map[string]interface{}{"role": "system", "content": "Hi can you tell me what is the factorial of 10?"}},
            "max_tokens": 50,
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

    // Extract the content from the JSON response
    content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
    fmt.Println(content)

}