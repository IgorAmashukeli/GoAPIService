package api

import (
	"HighTalent/internal/db"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllQuestionSuccess(t *testing.T) {
	ReadAllQuestionsFunc = func(_ gorm.Interface[db.Question], _ context.Context) ([]db.Question, error) {
		return []db.Question{
			{ID: 1, Text: "Test Question 1"},
			{ID: 2, Text: "Test Question 2"},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler := GetAllQuestionsHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Test Question 1")
	assert.Contains(t, w.Body.String(), "Test Question 2")
}

func TestGetAllQuestionsDbError(t *testing.T) {
	ReadAllQuestionsFunc = func(_ gorm.Interface[db.Question], _ context.Context) ([]db.Question, error) {
		return nil, errors.New("db failure")
	}

	req := httptest.NewRequest(http.MethodGet, "/questions", nil)
	w := httptest.NewRecorder()

	handler := GetAllQuestionsHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to fetch questions")
}

func TestCreateQuestionSuccess(t *testing.T) {
	CreateQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, text string) error {
		return nil
	}

	type RequestBody struct {
		Text string `json:"text"`
	}

	body, _ := json.Marshal(RequestBody{
		Text: "What is Golang?",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Question created")
}

func TestCreateQuestionInvalidJson(t *testing.T) {
	CreateQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, text string) error {
		return nil
	}

	body, _ := json.Marshal(map[string]string{
		"text":       "What is Golang?",
		"unexpected": "this should fail",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid JSON")
}

func TestCreateQuestionDbError(t *testing.T) {
	CreateQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, text string) error {
		return errors.New("db failure")
	}

	type RequestBody struct {
		Text string `json:"text"`
	}

	body, _ := json.Marshal(RequestBody{
		Text: "What is Golang?",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to create question")
}

func TestReadQuestionDbSuccess(t *testing.T) {
	ReadQuestionGetAllAnsFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, id uint) (db.Question, []db.Answer, error) {
		return db.Question{ID: 1, Text: "Test Question 1"},
			[]db.Answer{{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor"},
				{ID: 2, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Artem"}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	w := httptest.NewRecorder()

	handler := ReadQuestionHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Test Question 1")
}

func TestReadQuestionDbInvalidIdNotInt(t *testing.T) {
	ReadQuestionGetAllAnsFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, id uint) (db.Question, []db.Answer, error) {
		return db.Question{ID: 1, Text: "Test Question 1"},
			[]db.Answer{{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor"},
				{ID: 2, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Artem"}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/questions/abc", nil)
	w := httptest.NewRecorder()

	handler := ReadQuestionHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestReadQuestionDbInvalidIdNegative(t *testing.T) {
	ReadQuestionGetAllAnsFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, id uint) (db.Question, []db.Answer, error) {
		return db.Question{ID: 1, Text: "Test Question 1"},
			[]db.Answer{{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Text Answer 1"},
				{ID: 2, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Artem", Text: "Text Answer 1"}}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/questions/-1", nil)
	w := httptest.NewRecorder()

	handler := ReadQuestionHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestReadQuestionDbErrorInDb(t *testing.T) {
	ReadQuestionGetAllAnsFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, id uint) (db.Question, []db.Answer, error) {
		return db.Question{ID: 1, Text: "Test Question 1"},
			[]db.Answer{{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Text Answer 1"},
				{ID: 2, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Artem", Text: "Text Answer 1"}}, errors.New("db failure")
	}

	req := httptest.NewRequest(http.MethodGet, "/questions/1", nil)
	w := httptest.NewRecorder()

	handler := ReadQuestionHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to find a question")

}

func TestDeleteQuestionDbSuccess(t *testing.T) {
	ReadQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) (db.Question, error) {
		return db.Question{ID: 1, Text: "Test Question 1"}, nil
	}

	DeleteQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Question deleted")
}

func TestDeleteQuestionDbInvalidIdNotInt(t *testing.T) {
	ReadQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) (db.Question, error) {
		return db.Question{ID: 1, Text: "Test Question 1"}, nil
	}

	DeleteQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/questions/abc", nil)
	w := httptest.NewRecorder()

	handler := DeleteQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestDeleteQuestionDbInvalidIdNegative(t *testing.T) {
	ReadQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) (db.Question, error) {
		return db.Question{ID: 1, Text: "Test Question 1"}, nil
	}

	DeleteQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/questions/-1", nil)
	w := httptest.NewRecorder()

	handler := DeleteQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestDeleteQuestionDbNotFound(t *testing.T) {
	ReadQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) (db.Question, error) {
		return db.Question{ID: 1, Text: "Test Question 1"}, errors.New("db failure")
	}

	DeleteQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to find a question")
}

func TestDeleteQuestionDbErrorDb(t *testing.T) {
	ReadQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) (db.Question, error) {
		return db.Question{ID: 1, Text: "Test Question 1"}, nil
	}

	DeleteQuestionFunc = func(_ gorm.Interface[db.Question], _ context.Context, id uint) error {
		return errors.New("db failure")
	}

	req := httptest.NewRequest(http.MethodDelete, "/questions/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteQuestionHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to delete a question")
}

func TestCreateAnswerDbSuccess(t *testing.T) {
	CreateAnswerFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, _ uint, _ string, _ string) error {
		return nil
	}

	type RequestBody struct {
		Text   string `json:"text"`
		UserId string `json:"user_id"`
	}

	body, _ := json.Marshal(RequestBody{
		Text:   "What is Golang?",
		UserId: "Igor",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Answer created")

}

func TestCreateAnswerDbInvalidNotIntId(t *testing.T) {
	CreateAnswerFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, _ uint, _ string, _ string) error {
		return nil
	}

	type RequestBody struct {
		Text   string `json:"text"`
		UserId string `json:"user_id"`
	}

	body, _ := json.Marshal(RequestBody{
		Text:   "What is Golang?",
		UserId: "Igor",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions/abc/answers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestCreateAnswerDbInvalidNegativeId(t *testing.T) {
	CreateAnswerFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, _ uint, _ string, _ string) error {
		return nil
	}

	type RequestBody struct {
		Text   string `json:"text"`
		UserId string `json:"user_id"`
	}

	body, _ := json.Marshal(RequestBody{
		Text:   "What is Golang?",
		UserId: "Igor",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions/-1/answers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")
}

func TestCreateAnswerDbInvalidJson(t *testing.T) {
	CreateAnswerFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, _ uint, _ string, _ string) error {
		return nil
	}

	body, _ := json.Marshal(map[string]string{
		"text":       "What is Golang?",
		"user_id":    "Igor",
		"unexpected": "this should fail",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid JSON")
}

func TestCreateAnswerDB(t *testing.T) {
	CreateAnswerFunc = func(_ gorm.Interface[db.Question], _ gorm.Interface[db.Answer], _ context.Context, _ uint, _ string, _ string) error {
		return errors.New("db failure")
	}

	type RequestBody struct {
		Text   string `json:"text"`
		UserId string `json:"user_id"`
	}

	body, _ := json.Marshal(RequestBody{
		Text:   "What is Golang?",
		UserId: "Igor",
	})

	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := CreateAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to create answer")

}

func TestReadAnswerDB(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/answers/1", nil)
	w := httptest.NewRecorder()

	handler := ReadAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Test Answer 1")
}

func TestReadAnswerInvalidIdNotInt(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/answers/abc", nil)
	w := httptest.NewRecorder()

	handler := ReadAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")
}

func TestReadAnswerInvalidIdNegative(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/answers/-1", nil)
	w := httptest.NewRecorder()

	handler := ReadAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")
}

func TestReadAnswerBdError(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, errors.New("db failure")
	}

	req := httptest.NewRequest(http.MethodGet, "/answers/1", nil)
	w := httptest.NewRecorder()

	handler := ReadAnswerHandler(nil, nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to find answer")
}

func TestDeleteAnswerSuccess(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	DeleteAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/answers/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteAnswerHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "Answer deleted")

}

func TestDeleteAnswerInvalidIdNotInt(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	DeleteAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/answers/abc", nil)
	w := httptest.NewRecorder()

	handler := DeleteAnswerHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestDeleteAnswerInvalidIdNegative(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	DeleteAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/answers/-1", nil)
	w := httptest.NewRecorder()

	handler := DeleteAnswerHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "invalid id")

}

func TestDeleteAnswerNotFound(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, errors.New("db failure")
	}

	DeleteAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/answers/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteAnswerHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to find answer")
}

func TestDeleteAnswerErrorDb(t *testing.T) {
	ReadAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) (db.Answer, error) {
		return db.Answer{ID: 1, QuestionID: 1, Question: db.Question{ID: 1, Text: "Test Question 1"}, UserID: "Igor", Text: "Test Answer 1"}, nil
	}

	DeleteAnswerFunc = func(_ gorm.Interface[db.Answer], _ context.Context, _ uint) error {
		return errors.New("db failure")
	}

	req := httptest.NewRequest(http.MethodDelete, "/answers/1", nil)
	w := httptest.NewRecorder()

	handler := DeleteAnswerHandler(nil)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "failed to delete answer")
}
