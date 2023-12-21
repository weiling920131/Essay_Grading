package main

import (
    // "bufio"
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
    // file, err := os.Open("essay.txt") // Replace with the actual path to your .txt file
    // if err != nil {
    //     log.Fatalf("Error opening file: %v", err)
    // }
    // defer file.Close()

    // // Read the file content
    // scanner := bufio.NewScanner(file)
    // essayContent := ""
    // for scanner.Scan() {
    //     essayContent += scanner.Text() + "\n"
    // }

    // // Check for errors during scanning
    // if err := scanner.Err(); err != nil {
    //     log.Fatalf("Error reading from file: %v", err)
    // }

    // Set up the resty client and make the API call
    client := resty.New()
    
    // 輸入題目
    var topics = "With the popularity of social media, emojis have gradually become an indispensable form of expression in our lives. Although they can effectively save time and express emotions, it is important to use them appropriately in communication. Excessive reliance on emojis can lead to misunderstandings and a lack of genuineness in conveying emotions. Therefore, it is essential to maintain a balance in using emojis, combining them with appropriate text to better express our thoughts and feelings."
    // fmt.Printf("Essay Topics in English: ")
    // fmt.Scan(&topics)
    // 輸入文章
    var essay = "In this time and age, with the social media taking the world by storm, “emoji” gradually becomes an indispensable part of our everyday lives. When it comes to people’s love for emoji, there are several possible reasons that contribute to the phenomenon. First and foremost, people nowadays tend to lead a life packed with heavy workloads. Therefore, with emoji at hand, people are capable of delivering their feelings and thoughts without spending a great amount of time. Furthermore, emoji is simple but vivid. It makes it possible for the users to directly convey their messages with an intriguing facial expression. For instance, when elation and ecstasy run through my veins, Picture 1 would function greatly to visualize my feelings. Moreover, when I have a conflict with my friends with waves of indignation and wrath attacking me one after another, Picture 2 would become the best choice to display what I feel. On a whole, the use of emoji brings convenience to our lives and people may therefore consider it flawless when communicating, however, there is more to it than meets the eyes. \nAs a matter of fact, emoji tends to give rise to misunderstandings and confusions in many cases. I had a pertinent experience in the past. It was the day that I partaked in a speech contest in which I got a good grade after months of training and toil. The moment I got informed of the fascinating news, Picture 3 was what I delivered immediately to my friend. Nevertheless, he didn't response to me. He told me on the subsequent day that when receiving the smiling face with tears welling in the eyes, he could not determine whether I was choked to tears in melancholy or smiling with tears. In order to wrestle with the problem, the following are the solutions. First, it is advised to attach some explanatory words to the emoji to clearly shed light on the meaning. Second, if the receiver has mistaken the meaning in the past many times, words might be a better option. Last but not least, it is essential that we should think twice before sending them and ensure that the facial expression can match our messages well. In conclusion, if we take the aforementioned preventive measures, it is possible to convey our moods with emoji without false interpretations."
    // fmt.Printf("Essay Content in English: ")
    // fmt.Scan(&essay)

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
                        "content": "Please grade the following essay based on content, organization, grammar and sentence structure, and vocabulary and spelling. Each category should be scored from 1 to 5, with a total holistic score out of 20. Deduct 1 point in total score for essays significantly under the word count or not divided into paragraphs. Provide detailed scores and feedback for each category. By the way, the score should be consistent if the essay is the same.",
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
