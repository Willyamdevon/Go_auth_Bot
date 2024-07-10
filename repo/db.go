package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

const (
	hashTable = "tg_hashs"
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	//fmt.Println(fmt.Sprintf("%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.DBName))
	//db, err := sqlx.Connect("postgres", fmt.Sprintf("%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.DBName))
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tg_hashs (tg_id text, id_hash text, chat_id text, time timestamp)")
	if err != nil {
		panic(err)
	}

	return db, nil
}
