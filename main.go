package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)


func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Retrieve assistant ID from environment variables
    asstID := os.Getenv("GPT_ASST_ID")
    if asstID == "" {
        log.Fatal("Engine ID is not set in environment variables")
    }
    
    essayContent := "In this time and age, with the social media taking the world by storm, “emoji” gradually becomes an indispensable part of our everyday lives. When it comes to people’s love for emoji, there are several possible reasons that contribute to the phenomenon. First and foremost, people nowadays tend to lead a life packed with heavy workloads. Therefore, with emoji at hand, people are capable of delivering their feelings and thoughts without spending a great amount of time. Furthermore, emoji is simple but vivid. It makes it possible for the users to directly convey their messages with an intriguing facial expression. For instance, when elation and ecstasy run through my veins, Picture 1 would function greatly to visualize my feelings. Moreover, when I have a conflict with my friends with waves of indignation and wrath attacking me one after another, Picture 2 would become the best choice to display what I feel. On a whole, the use of emoji brings convenience to our lives and people may therefore consider it flawless when communicating, however, there is more to it than meets the eyes. As a matter of fact, emoji tends to give rise to misunderstandings and confusions in many cases. I had a pertinent experience in the past. It was the day that I partaked in a speech contest in which I got a good grade after months of training and toil. The moment I got informed of the fascinating news, Picture 3 was what I delivered immediately to my friend. Nevertheless, he didn’t response to me. He told me on the subsequent day that when receiving the smiling face with tears welling in the eyes, he could not determine whether I was choked to tears in melancholy or smiling with tears. In order to wrestle with the problem, the following are the solutions. First, it is advised to attach some explanatory words to the emoji to clearly shed light on the meaning. Second, if the receiver has mistaken the meaning in the past many times, words might be a better option. Last but not least, it is essential that we should think twice before sending them and ensure that the facial expression can match our messages well. In conclusion, if we take the aforementioned preventive measures, it is possible to convey our moods with emoji without false interpretations."

    // Create a thread
    thread, err := CreateThread()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Thread ID: %s\n", thread.ID)

    // Initialize the message (for test)
    message := essayContent

    for {
        // Create a message
        _, err = CreateThreadMessage(thread.ID, message)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        // Run the assistant
        _, err = CreateThreadRun(thread.ID)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        
        // Get the messages
        messages := GetMessages(thread.ID)
        fmt.Println("")
        fmt.Println("Messages:")

        // inverse order to get oldest messages first
        for i := len(messages.Data) - 1; i >= 0; i-- {
            message := messages.Data[i]
            fmt.Printf("Role: %s\n", message.Role)
            for _, content := range message.Content {
                fmt.Printf("Content Value: %s\n", content.Text.Value)
            }
        }

        // Read user input
        fmt.Print("Ask somethine: ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        message = scanner.Text()
    }
}
