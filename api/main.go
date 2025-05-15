package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"quiz/api/quizModel"
	"sync"

	"github.com/gorilla/mux"
)

var scoreHistory []int
var scoreLock sync.Mutex

func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizModel.Questions)
}

func submitAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var answers []quizModel.Answers
	if err := json.Unmarshal(body, &answers); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// calculate correct answers
	total := len(quizModel.CorrectAnswers)
	correct := 0
	for _, ans := range answers {
		if correctIndex, exists := quizModel.CorrectAnswers[ans.Id]; exists && correctIndex == ans.Choice-1 {
			correct++
		}
	}

	// Lock the score history to prevent two instances to change it at the same time
	scoreLock.Lock()
	scoreHistory = append(scoreHistory, correct)

	// Calculate the percentile
	percentile := calculatePercentile(correct, scoreHistory)
	scoreLock.Unlock()

	json.NewEncoder(w).Encode(quizModel.ScoreResponse{
		Score:      correct,
		Total:      total,
		Percentile: percentile,
	})
}

// Calculate the percentile based on the current score and all previous scores
func calculatePercentile(currentScore int, allScores []int) int {
	count := 0
	for _, s := range allScores {
		if s <= currentScore {
			count++
		}
	}
	percentile := int(math.Round(float64(count) / float64(len(allScores)) * 100))
	return percentile
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/questions", getQuestions).Methods("GET")
	router.HandleFunc("/answers", submitAnswers).Methods("POST")
	fmt.Println("Server started on port 8000")
	http.ListenAndServe(":8000", router)
}
