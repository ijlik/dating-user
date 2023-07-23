package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

// please ensure before using this fund
// error response must be on same variable
func txAction(tx *sqlx.Tx, err *error) {
	if *err != nil {
		log.Println(*err)
		errTx := tx.Rollback()
		if errTx != nil {
			log.Println("FAILED TO ROLLBACK: ", errTx)
		}
		return
	}

	errTx := tx.Commit()
	if errTx != nil {
		log.Println("FAILED TO COMMIT: ", errTx)
	}
}

func checkTagInt(tag sql.Result, msg string) error {
	tagInt, err := tag.RowsAffected()
	if err != nil {
		return err
	}

	if tagInt == 0 {
		return fmt.Errorf("failed to %s", msg)
	}

	return nil
}

type repo struct {
	conn *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepository {
	return &repo{db}
}
