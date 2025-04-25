package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"os"
)

const MySQLDeadlockErrorCode = 1213

func finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}
func NewMySQLConnection() (*sqlx.DB, error) {
	config := mysql.NewConfig()

	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_ADDR")
	config.User = os.Getenv("MYSQL_USER")
	config.Passwd = os.Getenv("MYSQL_PASSWORD")
	config.DBName = os.Getenv("MYSQL_DATABASE")
	config.ParseTime = true // with that parameter, we can use time.Time in mysqlHour.Hour

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to MySQL")
	}

	return db, nil
}
