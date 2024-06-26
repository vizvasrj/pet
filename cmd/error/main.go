package main

import (
	"database/sql"
	"fmt"
	"log"
	"src/myerror"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Aaaa struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func SomeFn(log *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, myerror.WrapError(log, err, "")
	}
	return db, nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	db, err := SomeFn(logger)
	if err != nil {
		log.Fatal(err)
	}
	aaaa := Aaaa{
		Db:     db,
		Logger: logger,
	}
	result, err := CheckOne(&aaaa)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

}

func (a *Aaaa) SelectOne() (int, error) {
	query := "SELECT * from some"
	one := 0
	err := a.Db.QueryRow(query).Scan(&one)
	if err != nil {
		return 0, myerror.WrapError(a.Logger, err, "")
	}
	if one != 1 {
		return one, myerror.New(a.Logger, "one is not 1")
	}
	return one, nil
}

func CheckOne(a *Aaaa) (string, error) {
	value, err := a.SelectOne()
	if err != nil {
		return fmt.Sprintf("value:: [%d]", value), myerror.WrapError(a.Logger, err, "")
	}
	return fmt.Sprintf("value:: [%d]", value), nil
}
