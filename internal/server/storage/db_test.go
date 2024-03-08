package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initDB() (*sql.DB, *sql.Tx) {
	DBDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `metrics_user`, `metrics_pass`, `metrics_db`)

	dbConnection, err := sql.Open("pgx", DBDsn)
	if err != nil {
		panic(err)
	}

	transaction, err := dbConnection.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return dbConnection, transaction
}

func TestDbExists(t *testing.T) {
	t.Skip("only manual use because depends on host")

	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	err := dbConnection.Ping()
	assert.NoError(t, err)

}

func TestDbTableExists(t *testing.T) {
	t.Skip("only manual use because depends on host")

	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	rawSQLString := "select * from metrics limit 1;"
	_, err := tx.Exec(rawSQLString)
	assert.NoError(t, err)
}

func TestCanMakeQuery(t *testing.T) {
	t.Skip("only manual use because depends on host")

	dbConnection, tx := initDB()
	defer dbConnection.Close()
	defer tx.Rollback()

	_, err := tx.Exec("insert into metrics values ('cnt', 'counter', null, 1)")
	assert.NoError(t, err)

	rawSQLString := "select count(*) from metrics;"
	rows, err := tx.Query(rawSQLString)
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())
}
