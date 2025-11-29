package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateAndReadQuestion(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	qdb, _, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateQuestion(qdb, ctx, "Hello?")
	assert.NoError(t, err)

	qs, err := ReadAllQuestions(qdb, ctx)
	assert.NoError(t, err)
	assert.Len(t, qs, 1)

	q, err := ReadQuestion(qdb, ctx, qs[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "Hello?", q.Text)

}

func TestDeleteQuestion_CascadeAnswers(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	qdb, adb, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateQuestion(qdb, ctx, "What is Go?")
	assert.NoError(t, err)
	qs, err := ReadAllQuestions(qdb, ctx)
	assert.NoError(t, err)
	qid := qs[0].ID
	err = CreateAnswer(adb, ctx, qid, "user001", "Go is awesome")
	assert.NoError(t, err)

	// Delete question → cascade answers
	err = DeleteQuestion(qdb, ctx, qid)
	assert.NoError(t, err)

	// Verify no answers left
	answers, err := adb.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, answers, 0)
}

func TestCreateAndReadAnswer(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	qdb, adb, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateQuestion(qdb, ctx, "2+2?")
	assert.NoError(t, err)
	qs, err := ReadAllQuestions(qdb, ctx)
	assert.NoError(t, err)
	qid := qs[0].ID

	err = CreateAnswer(adb, ctx, qid, "user123", "4")
	assert.NoError(t, err)

	as, err := adb.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, as, 1)

	a, err := ReadAnswer(adb, ctx, as[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "4", a.Text)
	assert.Equal(t, "user123", a.UserID)
}

func TestDeleteAnswer(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	qdb, adb, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateQuestion(qdb, ctx, "Delete test?")
	assert.NoError(t, err)
	qs, err := ReadAllQuestions(qdb, ctx)
	assert.NoError(t, err)
	qid := qs[0].ID
	err = CreateAnswer(adb, ctx, qid, "u1", "yes")
	assert.NoError(t, err)
	as, err := adb.Find(ctx)
	assert.NoError(t, err)

	err = DeleteAnswer(adb, ctx, as[0].ID)
	assert.NoError(t, err)

	as, err = adb.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, as, 0)
}

func TestCannotCreateAnswerForNonExistingQuestion(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	_, adb, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateAnswer(adb, ctx, 9999, "user123", "some text")
	assert.Error(t, err)
}

func TestUserCanCreateMultipleAnswers(t *testing.T) {
	database := setupTestDB(t)

	defer ResetInit(database, t)
	defer CloseDB(database)

	qdb, adb, ctx := ResetInitPrepare(database, t)

	assert.NotEqual(t, ctx, nil)

	err := CreateQuestion(qdb, ctx, "What is Go?")
	assert.NoError(t, err)
	qs, err := ReadAllQuestions(qdb, ctx)
	assert.NoError(t, err)
	qid := qs[0].ID

	err = CreateAnswer(adb, ctx, qid, "userA", "answer 1")
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	err = CreateAnswer(adb, ctx, qid, "userA", "answer 2")
	assert.NoError(t, err)

	as, err := adb.Find(ctx)
	assert.NoError(t, err)
	assert.Len(t, as, 2)
}

func ResetDatabase(database *sql.DB, t *testing.T) error {
	_, err := database.Exec(`
        DROP SCHEMA public CASCADE;
        CREATE SCHEMA public;
    `)
	if err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}
	return nil
}

func setupTestDB(t *testing.T) *sql.DB {
	database, connect_err := GetDatabase()
	if connect_err != nil {
		t.Fatalf("Failed to connect to test DB: %v", connect_err)
	}
	t.Cleanup(func() { CloseDB(database) })
	return database
}

func InitForTests(database *sql.DB) error {
	originalOutput := log.Default().Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(originalOutput)

	InitMigrationDir()
	return MigrateUp(database)
}

func ResetInit(database *sql.DB, t *testing.T) error {
	if err := ResetDatabase(database, t); err != nil {
		return fmt.Errorf("reset failed: %w", err)
	}
	if err := InitForTests(database); err != nil {
		return fmt.Errorf("init failed: %w", err)
	}
	return nil
}

func ResetInitPrepare(database *sql.DB, t *testing.T) (gorm.Interface[Question], gorm.Interface[Answer], context.Context) {
	if err := ResetInit(database, t); err != nil {
		return nil, nil, nil
	}
	return PrepareDBClients(database)
}
