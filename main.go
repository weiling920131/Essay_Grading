package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"

	// "github.com/fogleman/gg"
	// "github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
    "github.com/go-resty/resty/v2"
    "github.com/joho/godotenv"
)

const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func plotGrades(grades []int) {
	p := plot.New()

	// Create a bar chart for the grades
	bars, err := plotter.NewBarChart(plotter.Values{
		float64(grades[0]),
		float64(grades[1]),
		float64(grades[2]),
		float64(grades[3]),
	},
		vg.Points(40),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add the bars to the plot
	p.Add(bars)

	// Set labels for each bar
	p.NominalX("Organization", "Content", "Grammar", "Vocabulary")

	// Set the title and labels for the axes
	p.Title.Text = "Grades Distribution"
	p.X.Label.Text = "Categories"
	p.Y.Label.Text = "Grades"

    p.Y.Max = 5
	// Save the plot to a file (you can also use p.Show() to display the plot)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "grades_plot.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Grades plot saved to grades_plot.png")
}

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
                    "content": "Please grade the following essay based on content, organization, grammar and sentence structure, and vocabulary and spelling. Each category should be scored from 1 to 5 and socres should be integer number, with a total holistic score out of 20. Deduct 1 point in total score for essays significantly under the word count or not divided into paragraphs. Provide detailed scores and feedback for each category. Write the grade in the format Content@Organization@Grammer@Vocabulary@Total in the first line of your response. For example, the first line should be 4@5@4@5@18",
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
        // fmt.Println("Graded Essay Content:", content)
        lines := strings.Split(content, "\n")
        grade := strings.Split(lines[0], "@")
        var gradesInt []int
        for _, gradeStr := range grade {
            grade, err := strconv.Atoi(gradeStr)
            if err != nil {
                fmt.Printf("Error converting grade %s to integer: %v\n", gradeStr, err)
                return
            }
            gradesInt = append(gradesInt, grade)
        }
        fmt.Printf("Grade Organizatoin: %d\n", gradesInt[0]);
        fmt.Printf("Grade Content: %d\n", gradesInt[1]);
        fmt.Printf("Grade Grammar: %d\n", gradesInt[2]);
        fmt.Printf("Grade Vocabulary: %d\n", gradesInt[3]);
        for _, line := range lines[1:] {
            fmt.Printf("Feedback: %s\n", line)
        }
        plotGrades(gradesInt)
    } else {
        fmt.Println("No content returned in the response.")
    }
}
