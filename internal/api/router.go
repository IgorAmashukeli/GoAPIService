package api

import (
	"HighTalent/internal/db"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CreateApi(question_db gorm.Interface[db.Question], answer_db gorm.Interface[db.Answer]) *http.ServeMux {
	// Create a new ServeMux
	mux := http.NewServeMux()

	mux.HandleFunc("GET /questions/{id}", ReadQuestionHandler(question_db, answer_db))
	mux.HandleFunc("GET /questions/", GetAllQuestionsHandler(question_db))
	mux.HandleFunc("POST /questions/", CreateQuestionHandler(question_db))
	mux.HandleFunc("DELETE /questions/{id}", DeleteQuestionHandler(question_db))
	mux.HandleFunc("POST /questions/{id}/answers/", CreateAnswerHandler(question_db, answer_db))
	mux.HandleFunc("GET /answers/{id}", ReadAnswerHandler(question_db, answer_db))
	mux.HandleFunc("DELETE /answers/{id}", DeleteAnswerHandler(answer_db))

	fmt.Println("Server starting on port 8080...")

	return mux

}
