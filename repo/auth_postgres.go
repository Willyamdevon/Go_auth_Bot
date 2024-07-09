package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func CreateId(tgId int64, idHash string, db *sqlx.DB) (string, error) {
	var count int

	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id='%s'"), hashTable, tgId).Scan(&count)
	switch {
	case err != nil:
		return "", err
	default:
		if count == 0 {
			var id string
			query := fmt.Sprintf("INSERT INTO %s (tg_id, id_hash, time) VALUES ($1, $2, $3) RETURNING id_hash", hashTable)

			row := db.QueryRow(query, tgId, idHash, time.Now())
			if err := row.Scan(&id); err != nil {
				return "", err
			}

			return id, nil
		}
		return "Уже есть", nil
	}

}

func GetCurenHash(tgId int64, db *sqlx.DB) (string, error) {
	var hash string

	err := db.QueryRow(fmt.Sprintf("SELECT hash FROM %s WHERE id='%s'"), hashTable, tgId).Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil

}
