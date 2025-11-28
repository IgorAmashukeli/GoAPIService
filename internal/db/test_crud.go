package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func TestCreateAndReadQuestion(database *sql.DB, qdb gorm.Interface[Question],
	ctx context.Context) {

	err := CreateQuestion(qdb, ctx, "Hello?")
	if err != nil {
		CloseDB(database)
		log.Fatalf("CreateQuestion failed: %v", err)
	}

	qs, err := ReadAllQuestions(qdb, ctx)
	if err != nil {
		CloseDB(database)
		log.Fatalf("ReadAllQuestions failed: %v", err)
	}
	if len(qs) != 1 {
		CloseDB(database)
		log.Fatalf("expected 1 question, got %d", len(qs))
	}

	q, err := ReadQuestion(qdb, ctx, qs[0].ID)
	if err != nil {
		CloseDB(database)
		log.Fatalf("ReadQuestion failed: %v", err)
	}
	if q.Text != "Hello?" {
		CloseDB(database)
		log.Fatalf("expected text 'Hello?', got '%s'", q.Text)
	}

	log.Println("Succesfull read and create of questions")
}

func TestDeleteQuestion_CascadeAnswers(database *sql.DB, qdb gorm.Interface[Question],
	adb gorm.Interface[Answer],
	ctx context.Context) {

	err := CreateQuestion(qdb, ctx, "What is Go?")
	if err != nil {
		CloseDB(database)
		log.Fatalf("CreateQuestion failed: %v", err)
	}

	qs, _ := ReadAllQuestions(qdb, ctx)
	qid := qs[0].ID

	err = CreateAnswer(adb, ctx, qid, "user001", "Go is awesome")
	if err != nil {
		CloseDB(database)
		log.Fatalf("CreateAnswer failed: %v", err)
	}

	err = DeleteQuestion(qdb, ctx, qid)
	if err != nil {
		CloseDB(database)
		log.Fatalf("DeleteQuestion failed: %v", err)
	}

	var answers []Answer
	answers, err = adb.Find(ctx)
	if err != nil {
		CloseDB(database)
		log.Fatalf(" reading answers after delete: %v", err)
	}
	if len(answers) != 0 {
		CloseDB(database)
		log.Fatalf("expected 0 answers after cascade delete, got %d", len(answers))
	}

	log.Println("Succesfull delete of questions and cascade delete of answers")
}

func TestCreateAndReadAnswer(database *sql.DB, qdb gorm.Interface[Question],
	adb gorm.Interface[Answer],
	ctx context.Context) {

	if err := CreateQuestion(qdb, ctx, "2+2?"); err != nil {
		CloseDB(database)
		log.Fatalf("CreateQuestion failed: %v", err)
	}
	qs, _ := ReadAllQuestions(qdb, ctx)
	qid := qs[0].ID

	if err := CreateAnswer(adb, ctx, qid, "user123", "4"); err != nil {
		CloseDB(database)
		log.Fatalf("CreateAnswer failed: %v", err)
	}

	as, err := adb.Find(ctx)
	if err != nil {
		CloseDB(database)
		log.Fatalf("Find failed: %v", err)
	}
	if len(as) != 1 {
		CloseDB(database)
		log.Fatalf("expected 1 answer, got %d", len(as))
	}

	a, err := ReadAnswer(adb, ctx, as[0].ID)
	if err != nil {
		CloseDB(database)
		log.Fatalf("ReadAnswer failed: %v", err)
	}
	if a.Text != "4" {
		CloseDB(database)
		log.Fatalf("expected '4', got '%s'", a.Text)
	}
	if a.UserID != "user123" {
		CloseDB(database)
		log.Fatalf("wrong user_id: '%s'", a.UserID)
	}

	log.Println("Succesfull read and create of answers")
}

func TestDeleteAnswer(database *sql.DB, qdb gorm.Interface[Question],
	adb gorm.Interface[Answer],
	ctx context.Context) {

	if err := CreateQuestion(qdb, ctx, "Delete test?"); err != nil {
		CloseDB(database)
		log.Fatalf("CreateQuestion failed: %v", err)
	}
	qs, _ := ReadAllQuestions(qdb, ctx)
	qid := qs[0].ID

	if err := CreateAnswer(adb, ctx, qid, "u1", "yes"); err != nil {
		CloseDB(database)
		log.Fatalf("CreateAnswer failed: %v", err)
	}

	as, _ := adb.Find(ctx)

	if err := DeleteAnswer(adb, ctx, as[0].ID); err != nil {
		CloseDB(database)
		log.Fatalf("DeleteAnswer failed: %v", err)
	}

	as, err := adb.Find(ctx)
	if err != nil {
		CloseDB(database)
		log.Fatalf("Find after delete failed: %v", err)
	}
	if len(as) != 0 {
		CloseDB(database)
		log.Fatalf("expected 0 answers, got %d", len(as))
	}

	log.Println("Succesfull delete of answers")
}

func TestCannotCreateAnswerForNonExistingQuestion(database *sql.DB,
	adb gorm.Interface[Answer],
	ctx context.Context) {

	// try to create answer for question_id = 9999
	err := CreateAnswer(adb, ctx, 9999, "user123", "some text")

	if err == nil {
		CloseDB(database)
		log.Fatalf("expected error creating answer for non-existing question, got nil")
	}

	log.Println("Succesfull test on creating answer for non existing question")

}

func TestUserCanCreateMultipleAnswers(database *sql.DB, qdb gorm.Interface[Question],
	adb gorm.Interface[Answer], ctx context.Context) {

	err := CreateQuestion(qdb, ctx, "What is Go?")
	if err != nil {
		CloseDB(database)
		log.Fatalf("failed to create question: %v", err)
	}

	qs, _ := ReadAllQuestions(qdb, ctx)
	if len(qs) != 1 {
		CloseDB(database)
		log.Fatalf("expected 1 question, got %d", len(qs))
	}
	qID := qs[0].ID

	err = CreateAnswer(adb, ctx, qID, "userA", "answer 1")
	if err != nil {
		CloseDB(database)
		log.Fatalf("failed to create answer #1: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	err = CreateAnswer(adb, ctx, qID, "userA", "answer 2")
	if err != nil {
		CloseDB(database)
		log.Fatalf("failed to create answer #2: %v", err)
	}

	// check stored answers
	as, err := adb.Find(ctx)
	if err != nil {
		CloseDB(database)
		log.Fatalf("failed reading answers: %v", err)
	}
	if len(as) != 2 {
		CloseDB(database)
		log.Fatalf("expected 2 answers from same user, got %d", len(as))
	}

	log.Println("Succesfull test on creating answer for non existing question")
}

func ResetDatabase(database *sql.DB) error {
	_, err := database.Exec(`
        DROP SCHEMA public CASCADE;
        CREATE SCHEMA public;
    `)

	if err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	return nil
}

func InitForTests(database *sql.DB) error {
	InitMigrationDir()
	if err := MigrateUp(database); err != nil {
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}

func ResetInit(database *sql.DB) {
	err_res := ResetDatabase(database)
	if err_res != nil {
		CloseDB(database)
		log.Fatalf("%v", err_res)

	}
	err_init := InitForTests(database)
	if err_init != nil {
		CloseDB(database)
		log.Fatalf("%v", err_init)
	}
}

func ResetInitPrepare(database *sql.DB) (gorm.Interface[Question],
	gorm.Interface[Answer], context.Context) {

	ResetInit(database)

	return PrepareDBClients(database)
}

func TestCRUD(database *sql.DB) {

	question_db, answer_db, ctx := ResetInitPrepare(database)
	TestCreateAndReadQuestion(database, question_db, ctx)

	question_db, answer_db, ctx = ResetInitPrepare(database)
	TestDeleteQuestion_CascadeAnswers(database, question_db, answer_db, ctx)

	question_db, answer_db, ctx = ResetInitPrepare(database)
	TestCreateAndReadAnswer(database, question_db, answer_db, ctx)

	question_db, answer_db, ctx = ResetInitPrepare(database)
	TestDeleteAnswer(database, question_db, answer_db, ctx)

	question_db, answer_db, ctx = ResetInitPrepare(database)
	TestCannotCreateAnswerForNonExistingQuestion(database, answer_db, ctx)

	question_db, answer_db, ctx = ResetInitPrepare(database)
	TestUserCanCreateMultipleAnswers(database, question_db, answer_db, ctx)

	ResetInit(database)
}
