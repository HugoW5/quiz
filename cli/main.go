package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"quiz/api/quizModel"
	"strings"

	"github.com/spf13/cobra"
)

type Question struct {
	Id      string
	Text    string
	Choiecs []string
}

type Answers struct {
	Id     string
	Choice int
}

func main() {
	fmt.Println("Välkommen till quizet, välj ett alternativ:")
	fmt.Println("1. Starta quizet")
	fmt.Println("2. Avsluta programmet")

	reader := bufio.NewReader(os.Stdin) // Consolebuffer reader

	for {
		fmt.Print("\nDitt val: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input) // remove whitespace

		switch input {
		case "1":
			fmt.Println("\nStartar quizet...")
			if err := startCommand.Execute(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "2":
			fmt.Println("\nAvslutar programmet...")
			return
		default:
			fmt.Println("\nOgiltigt val, försök igen.")
		}
	}
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Hämtar frågor från quiz-API och starta quizet",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get("http://localhost:8000/questions")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var questions []Question
		var answers []Answers = []Answers{}

		if err := json.NewDecoder(resp.Body).Decode(&questions); err != nil {
			return err
		}
		reader := bufio.NewReader(os.Stdin) // Consolebuffer reader

		for _, q := range questions {
			fmt.Printf("%s: %s\n", q.Id, q.Text)
			for i, a := range q.Choiecs {
				fmt.Printf("  %d. %s\n", i+1, a)
			}

		inputLoop:
			for {
				fmt.Print("\nDitt val: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input) // remove whitespace
				switch input {
				case "1":
					answers = append(answers, Answers{Id: q.Id, Choice: 1})
					break inputLoop
				case "2":
					answers = append(answers, Answers{Id: q.Id, Choice: 2})
					break inputLoop
				case "3":
					answers = append(answers, Answers{Id: q.Id, Choice: 3})
					break inputLoop
				default:
					fmt.Println("\nOgiltigt val, försök igen.")
				}
			}

		}
		submitAnswers(answers)
		return nil
	},
}

func submitAnswers(answers []Answers) error {
	// Convert the answers to JSON
	payload, err := json.Marshal(answers)
	if err != nil {
		return fmt.Errorf("could not marshal answers: %w", err)
	}

	// Send the POST request
	resp, err := http.Post("http://localhost:8000/answers", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	// Parse the response
	var scoreResponse quizModel.ScoreResponse
	if err := json.Unmarshal(body, &scoreResponse); err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}

	// Show result
	fmt.Printf("\nQuiz Resultat:\nScore: %d/%d\nDu var bättre än: %d%% av alla som tagit quizet\n", scoreResponse.Score, scoreResponse.Total, scoreResponse.Percentile)

	return nil
}
