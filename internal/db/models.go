package db

import (
	"context"
	"database/sql"
    "time"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

type Question struct {
    ID        uint      `gorm:"primaryKey"`
    Text      string    `gorm:"type:text"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Answer struct {
    ID         uint      `gorm:"primaryKey"`
    QuestionID uint      `gorm:"not null"` // FK column
    Question   Question  `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
    UserID     string    `gorm:"type:text"`
    Text       string    `gorm:"type:text"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func GetGormDatabase(sql_db *sql.DB) (*gorm.DB, error) {
	gorm_db, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sql_db,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })

	if err != nil {
		return nil, err
	}

	sqldbgorm, gorm_err := gorm_db.DB();
	if gorm_err != nil {
		return gorm_db, gorm_err
	}

	sqldberr := sqldbgorm.Ping()
	if sqldberr != nil {
		return gorm_db, sqldberr
	}

	return gorm_db, nil
}


func GetContext() context.Context {
  return context.Background()
}


func GetQuestionsDB(gorm_db *gorm.DB ) gorm.Interface[Question] {
	return gorm.G[Question](gorm_db)
}

func GetAnswersDB(gorm_db *gorm.DB) gorm.Interface[Answer] {
	return gorm.G[Answer](gorm_db)
}

