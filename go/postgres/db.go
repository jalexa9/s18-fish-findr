package postgres

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	// postgres driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DB embeds an SQLx database driver.
type DB struct {
	*sqlx.DB
	source string
}

// CreateDB connects to a Postgres database.
func CreateDB() (*DB, error) {
	dataSource := os.Getenv("DB_URL")
	conn, err := sqlx.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	db := DB{
		DB:     conn,
		source: dataSource,
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &db, nil
}

// Transact provides a wrapper for transactions.
func (db *DB) Transact(txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			default:
				err = fmt.Errorf("%s", r)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	err = txFunc(tx)
	return err
}

// Close closes all DB connections.
func (db *DB) Close() error {
	err := db.DB.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close db")
	}
	return nil
}

func buildValues(numCols int, count int) string {
	placeholders := make([]string, numCols)
	for c := 0; c < numCols; c++ {
		placeholders[c] = "?"
	}
	counts := make([]string, count)
	holder := fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
	for i := 0; i < count; i++ {
		counts[i] = holder
	}
	return sqlx.Rebind(sqlx.DOLLAR, strings.Join(counts, ","))
}

// CloseOrFail closes all DB connections, or fails the tests.
func (db *DB) CloseOrFail(t require.TestingT) {
	err := db.Close()
	require.Nil(t, err)
}

// Used in the testing of the DB.
func (db *DB) assertCountRows(t *testing.T, num int, table string) {
	var count int
	err := db.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", table))
	assert.Nil(t, err)
	assert.Equal(t, num, count)
}

// Used in the testing of the DB.
func (db *DB) assertValueOfWhere(t *testing.T, num int, table string, value string) {
	var count int
	err := db.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", table, value))
	assert.Nil(t, err)
	assert.Equal(t, num, count)
}

// CleanDB removes everything from the profile tables in the DB.
func (db *DB) CleanDB(t *testing.T) error {
	// TODO, add removes for more than just the profile table.
	query := `DELETE FROM user_interest`
	_, err := db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "Error deleting from user_interest table.")
	}
	query = `DELETE FROM profile`
	_, err = db.Exec(query)
	return errors.Wrapf(err, "Error deleting from profile table.")
}
