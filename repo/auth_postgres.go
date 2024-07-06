package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func CreateId(tgId int64, idHash string, db *sqlx.DB) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (tg_id, id_hash, time) VALUES ($1, $2. $3) RETURNING id_hash", hashTable)

	row := db.QueryRow(query, tgId, idHash, time.Now())
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
