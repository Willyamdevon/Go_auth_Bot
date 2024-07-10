package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func CreateId(tgId int64, idHash string, db *sqlx.DB) (string, error) {
	var count int

	count, err := CountOfID(tgId, db)
	switch {
	case err != nil:
		return "", err
	default:
		if count == 0 {
			var id string

			query := fmt.Sprintf("INSERT INTO %s (tg_id, id_hash, time) VALUES ($1, $2, $3) RETURNING id_hash", hashTable)

			row := db.QueryRow(query, tgId, idHash, time.Now().Add(time.Hour*12).UTC())
			if err := row.Scan(&id); err != nil {
				fmt.Println(err)
				return "", err
			}

			return id, nil
		}
		return "Уже есть", nil
	}

}

func GetCurentHash(tgId int64, db *sqlx.DB) (string, string, error) {
	var hash string
	var timeStart time.Time

	query := fmt.Sprintf("SELECT id_hash, time FROM %s WHERE tg_id=$1", hashTable)

	err := db.QueryRow(query, tgId).Scan(&hash, &timeStart)
	if err != nil {
		return "", "", err
	}

	currentTime := time.Now().UTC()
	duration := timeStart.Sub(currentTime)

	seconds := int(duration.Seconds())
	hours := seconds / 3600
	minutes := (seconds - hours*3600) / 60
	remainingSeconds := seconds - hours*3600 - minutes*60

	result := fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, remainingSeconds)

	return hash, result, nil
}

func CountOfID(tgId int64, db *sqlx.DB) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tg_id=$1", hashTable)

	if err := db.QueryRow(query, tgId).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}
