package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// findQuestionsChoicesAndAnswers finds questions, their respective choices, and answers.
func findQuestionsChoicesAndAnswers(lines []string) (map[string][]string, map[string]string, error) {
	reQuestion := regexp.MustCompile(`^[A-Z].*[:?]$`)
	reChoice := regexp.MustCompile(`^--.*$`)
	reAnswer := regexp.MustCompile(`^ans:[a-zA-Z]$`)

	questionsAndChoices := make(map[string][]string)
	questionAndAnswer := make(map[string]string)
	var currentQuestion string

	for _, line := range lines {
		if reQuestion.MatchString(line) {
			currentQuestion = line
		} else if reChoice.MatchString(line) {
			choice := line[2:] // Remove '--'
			if currentQuestion != "" {
				questionsAndChoices[currentQuestion] = append(questionsAndChoices[currentQuestion], choice)
			}
		} else if reAnswer.MatchString(line) {
			answer := line[4:] // Remove 'ans:'
			if currentQuestion != "" {
				questionAndAnswer[currentQuestion] = answer
			}
		}
	}

	return questionsAndChoices, questionAndAnswer, nil
}

func question(text string) (string, []string, string) {
	questionRegex := regexp.MustCompile(`[A-Z][^?|:]*[?:]`)
	questions := questionRegex.FindAllString(text, -1)

	choiceRegex := regexp.MustCompile(`--[^/n]*`)
	choices := choiceRegex.FindAllString(text, -1)

	answerRegex := regexp.MustCompile(`ans:[a-zA-Z]`)
	answer := answerRegex.FindString(text)

	if len(questions) > 0 && len(choices) > 0 {
		randIndex := rand.Intn(len(questions))
		selectedQuestion := questions[randIndex]

		return selectedQuestion, choices, answer
	}
	return "", nil, ""
}

func selectChoice(choices []string) string {
	fmt.Println("Choices:")
	for index, choice := range choices {
		fmt.Printf("%c.%s\n", 'a'+rune(index), choice)
	}

	fmt.Println("Please enter your choice (a, b, c, etc.):")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userChoice := strings.TrimSpace(userInput)
	return userChoice
}

func choice(choices []string) []string {
	labels := []string{"a", "b", "c", "d"}
	var labeledChoices []string

	for i, ch := range choices {
		labeledChoices = append(labeledChoices, fmt.Sprintf("%s.%s", labels[i], strings.TrimPrefix(ch, "--")))
	}
	return labeledChoices
}

func nextQuestion(questions []string) string {
	if len(questions) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(questions))
	return questions[randomIndex]
}

func handler(w http.ResponseWriter, r *http.Request) {
	text := `Unless otherwise posted, what is the basic speed limit outside a city, town or village on a primary highway?
--100 km/h
--90 km/h
--110 km/h
--80 km/h
ans:a
When a driver is stopping behind another vehicle in traffic:
--Immediately
--Within 5 seconds
--Within 2 seconds
ans:b`

	selectedQuestion, choices, _ := question(text)
	labeledChoices := choice(choices)

	fmt.Fprintf(w, "<h1>%s</h1>", selectedQuestion)
	for _, ch := range labeledChoices {
		fmt.Fprintf(w, "<p>%s</p>", ch)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open("french.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	questionsAndChoices, questionAndAnswer, err := findQuestionsChoicesAndAnswers(lines)
	if err != nil {
		fmt.Println("Error finding questions and choices:", err)
		return
	}

	var questions []string
	for q := range questionsAndChoices {
		questions = append(questions, q)
	}

	for {
		question := nextQuestion(questions)
		if question == "" {
			fmt.Println("No more questions.")
			break
		}
		fmt.Println("Question:", question)
		choices := questionsAndChoices[question]
		userChoice := selectChoice(choices)

		correctAnswer := questionAndAnswer[question]
		if userChoice == correctAnswer {
			fmt.Println("Correct!")
		} else {
			fmt.Println("Maybe the correct answer is:", correctAnswer)
		}

		fmt.Println("Would you like to answer the next question? (absolutely/later)")
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		if strings.TrimSpace(userInput) == "later" {
			break
		}
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Word in the host: %s", r.Host)
	// })

	// fmt.Println("Starting server at :8080")
	// http.ListenAndServe(":8080", nil)

	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":9090", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
