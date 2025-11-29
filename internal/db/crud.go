package db

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func PrepareDBClients(database *sql.DB) (gorm.Interface[Question],
	gorm.Interface[Answer], context.Context) {
	gorm_db, gorm_err := GetGormDatabase(database)
	if gorm_err != nil {
		CloseDB(database)
		return nil, nil, nil
	}

	question_db := GetQuestionsDB(gorm_db)
	answer_db := GetAnswersDB(gorm_db)

	ctx := GetContext()

	return question_db, answer_db, ctx
}

func CreateQuestion(question_db gorm.Interface[Question], ctx context.Context, text string) error {
	q := Question{Text: text, CreatedAt: time.Now()}
	err := question_db.Create(ctx, &q)
	return err
}

func ReadAllQuestions(question_db gorm.Interface[Question], ctx context.Context) ([]Question, error) {
	return question_db.Find(ctx)
}

func ReadQuestion(question_db gorm.Interface[Question], ctx context.Context, id uint) (Question, error) {
	q, err := question_db.Where("id = ?", id).First(ctx)
	return q, err
}

func ReadQuestionGetAllAnswers(question_db gorm.Interface[Question], answer_db gorm.Interface[Answer], ctx context.Context, id uint) (Question, []Answer, error) {
	q, err := ReadQuestion(question_db, ctx, id)
	if err != nil {
		return q, nil, err
	}

	answers, answer_err := answer_db.Where("question_id = ?", id).Find(ctx)
	if answer_err != nil {
		return q, nil, answer_err
	}

	return q, answers, nil
}

func DeleteQuestion(question_db gorm.Interface[Question], ctx context.Context, id uint) error {
	_, err := question_db.Where("id = ?", id).Delete(ctx)
	return err
}

func CreateAnswer(answer_db gorm.Interface[Answer], ctx context.Context, question_id uint, user_id string, text string) error {
	a := Answer{QuestionID: question_id, UserID: user_id, Text: text, CreatedAt: time.Now()}
	err := answer_db.Create(ctx, &a)
	return err
}

func CreateValidatedAnswer(question_db gorm.Interface[Question], answer_db gorm.Interface[Answer], ctx context.Context,
	question_id uint, user_id string, text string) error {

	_, err := ReadQuestion(question_db, ctx, question_id)

	if err != nil {
		return err
	}

	return CreateAnswer(answer_db, ctx, question_id, user_id, text)

}

func ReadAnswer(answer_db gorm.Interface[Answer], ctx context.Context, id uint) (Answer, error) {
	a, err := answer_db.Where("id = ?", id).First(ctx)
	return a, err
}

func DeleteAnswer(answer_db gorm.Interface[Answer], ctx context.Context, id uint) error {
	_, err := answer_db.Where("id = ?", id).Delete(ctx)
	return err
}
