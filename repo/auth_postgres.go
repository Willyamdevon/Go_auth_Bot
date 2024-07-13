package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func CreateId(tgId int64, idHash string, chatId int64, db *sqlx.DB) (string, error) {
	var count int

	count, err := CountOfID(tgId, db)
	switch {
	case err != nil:
		return "", err
	default:
		if count == 0 {
			var id string

			query := fmt.Sprintf("INSERT INTO %s (tg_id, id_hash, chat_id, time) VALUES ($1, $2, $3, $4) RETURNING id_hash", hashTable)
			// TODO: добавить в бд столбец-статус - ссылка уже использованна
			row := db.QueryRow(query, tgId, idHash, chatId, time.Now().Add(time.Hour*12).UTC())
			if err := row.Scan(&id); err != nil {
				fmt.Println(err)
				return "", err
			}

			return id, nil
		}
		return "Уже есть", nil
	}

}

func GetCurentLink(tgId int64, db *sqlx.DB) (string, string, error) {
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

	if seconds > 0 && hours > 0 && remainingSeconds > 0 {
		result := fmt.Sprintf("ссылка действует %d hours, %d minutes, %d seconds", hours, minutes, remainingSeconds)

		return hash, result, nil
	}

	return hash, "", nil
}

func CountOfID(tgId int64, db *sqlx.DB) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tg_id=$1", hashTable)

	if err := db.QueryRow(query, tgId).Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func DeleteLink(tgId int64, db *sqlx.DB) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE tg_id=$1", hashTable)
	if _, err := db.Exec(query, tgId); err != nil {
		return err
	}
	return nil
}

func GetCurentTime(tgId int64, db *sqlx.DB) (string, error) {
	count, err := CountOfID(tgId, db)
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", nil
	}

	var timeStart time.Time

	query := fmt.Sprintf("SELECT time FROM %s WHERE tg_id=$1", hashTable)

	row := db.QueryRow(query, tgId)
	if err := row.Scan(&timeStart); err != nil {
		return "", err
	}

	currentTime := time.Now().UTC()
	duration := timeStart.Sub(currentTime)

	seconds := int(duration.Seconds())
	hours := seconds / 3600
	minutes := (seconds - hours*3600) / 60
	remainingSeconds := seconds - hours*3600 - minutes*60

	if seconds > 0 && hours > 0 && remainingSeconds > 0 {
		result := fmt.Sprintf("ссылка действует %d hours, %d minutes, %d seconds", hours, minutes, remainingSeconds)

		return result, nil
	}

	return "", nil
}

func GetCurentHash(tgId int64, db *sqlx.DB) (string, error) {
	count, err := CountOfID(tgId, db)
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", nil
	}

	var hash string

	query := fmt.Sprintf("SELECT id_hash FROM %s WHERE tg_id=$1", hashTable)
	row := db.QueryRow(query, tgId)
	if err := row.Scan(&hash); err != nil {
		return "", err
	}
	log.Println(hash)
	return hash, nil
}
