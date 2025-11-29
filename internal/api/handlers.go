package api

import (
	"HighTalent/internal/db"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var ReadAllQuestionsFunc = db.ReadAllQuestions

func GetAllQuestionsHandler(question_db gorm.Interface[db.Question]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		questions, err := ReadAllQuestionsFunc(question_db, r.Context())
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return

		}

		w.Header().Set("Content-Type", "application/json")

		encode_err := json.NewEncoder(w).Encode(questions)

		if encode_err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

type CreateQuestionRequest struct {
	Text string `json:"text"`
}

var CreateQuestionFunc = db.CreateQuestion

func CreateQuestionHandler(question_db gorm.Interface[db.Question]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		var req CreateQuestionRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := CreateQuestionFunc(question_db, r.Context(), req.Text); err != nil {
			http.Error(w, "failed to create question", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Question created"})

	}
}

var ReadQuestionGetAllAnsFunc = db.ReadQuestionGetAllAnswers
var ReadQuestionFunc = db.ReadQuestion

type GetQAllAnsResponse struct {
	Question db.Question `json:"question"`
	Answers  []db.Answer `json:"answers"`
}

func ReadQuestionHandler(question_db gorm.Interface[db.Question], answer_db gorm.Interface[db.Answer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/questions/")
		int_id, err := strconv.Atoi(path)
		if err != nil || int_id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		question, answers, read_err := ReadQuestionGetAllAnsFunc(question_db, answer_db, r.Context(), uint(int_id))

		if read_err != nil {
			if errors.Is(read_err, gorm.ErrRecordNotFound) {
				http.Error(w, "question not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")

		value_answer := GetQAllAnsResponse{question, answers}

		encode_err := json.NewEncoder(w).Encode(value_answer)

		if encode_err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

var DeleteQuestionFunc = db.DeleteQuestion

func DeleteQuestionHandler(question_db gorm.Interface[db.Question]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/questions/")
		int_id, err := strconv.Atoi(path)
		if err != nil || int_id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		_, read_err := ReadQuestionFunc(question_db, r.Context(), uint(int_id))

		if read_err != nil {
			if errors.Is(read_err, gorm.ErrRecordNotFound) {
				http.Error(w, "question not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		delete_err := DeleteQuestionFunc(question_db, r.Context(), uint(int_id))

		if delete_err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Question deleted"})

	}
}

type CreateAnswerRequest struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

var CreateAnswerFunc = db.CreateValidatedAnswer

func CreateAnswerHandler(question_db gorm.Interface[db.Question], answer_db gorm.Interface[db.Answer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/questions/")
		path = strings.TrimSuffix(path, "/answers/")
		int_id, err := strconv.Atoi(path)
		if err != nil || int_id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		var req CreateAnswerRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		decoder_err := decoder.Decode(&req)
		if decoder_err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := CreateAnswerFunc(question_db, answer_db, r.Context(), uint(int_id), req.UserId, req.Text); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "question not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Answer created"})
	}
}

var ReadAnswerFunc = db.ReadAnswer

func ReadAnswerHandler(question_db gorm.Interface[db.Question], answer_db gorm.Interface[db.Answer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/answers/")
		int_id, err := strconv.Atoi(path)
		if err != nil || int_id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		answer, answer_err := ReadAnswerFunc(answer_db, r.Context(), uint(int_id))
		if answer_err != nil {
			if errors.Is(answer_err, gorm.ErrRecordNotFound) {
				http.Error(w, "answer not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")

		encode_err := json.NewEncoder(w).Encode(answer)

		if encode_err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

var DeleteAnswerFunc = db.DeleteAnswer

func DeleteAnswerHandler(answer_db gorm.Interface[db.Answer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/answers/")
		int_id, err := strconv.Atoi(path)
		if err != nil || int_id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		_, answer_err := ReadAnswerFunc(answer_db, r.Context(), uint(int_id))
		if answer_err != nil {
			if errors.Is(answer_err, gorm.ErrRecordNotFound) {
				http.Error(w, "answer not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")

		delete_answer_err := DeleteAnswerFunc(answer_db, r.Context(), uint(int_id))
		if delete_answer_err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Answer deleted"})

	}
}
