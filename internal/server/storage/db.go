package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

type DBStorage struct {
	DBConnection *sql.DB
}

func NewDBStorage(dsn string) *DBStorage {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		//
		logger.ZapSugarLogger.Panicln("could not connect to db using dsn: " + dsn)
	}

	createType := `DO $$ BEGIN    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metric_type') THEN
	       CREATE TYPE metric_type AS
	       ENUM(
	    		'gauge', 'counter'
	       );
	   END IF; END$$`
	_, err = conn.Exec(createType)
	if err != nil {
		panic(err)
	}

	createTable := `create table if not exists metrics
	(
	name varchar(255) primary key ,
	type metric_type not null default 'gauge',
	value double precision default null,
	delta numeric default null
	);`
	_, err = conn.Exec(createTable)
	if err != nil {
		panic(err)
	}

	return &DBStorage{
		DBConnection: conn,
	}
}

func (strg *DBStorage) Clear() {
	_, err := strg.DBConnection.Exec("delete from metrics")
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
	}
}

func (strg *DBStorage) AddCounter(ctx context.Context, key string, value int64) {
	query := `
			INSERT INTO metrics ( name, type, delta)
			VALUES ($1, $2, $3)
			ON CONFLICT(name) 
			DO UPDATE SET
			  delta = metrics.delta + $3;`

	_, err := strg.DBConnection.ExecContext(ctx, query, key, types.Counter, value)
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not add counter", err.Error())
	}
}

func (strg *DBStorage) SetGauge(ctx context.Context, key string, value float64) {
	query := `
			INSERT INTO metrics ( name, type, value)
			VALUES ($1, $2, $3)
			ON CONFLICT(name) 
			DO UPDATE SET
			  value = $3;`

	_, err := strg.DBConnection.ExecContext(ctx, query, key, types.Gauge, value)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
	}
}

func (strg *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	query := `select value from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(ctx, query, name, types.Gauge)
	if row.Err() != nil {
		return 0, fmt.Errorf(exceptions.UnknownMetricName + " " + row.Err().Error())
	}

	var gauge float64
	err := row.Scan(&gauge)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, fmt.Errorf(exceptions.UnknownMetricName)
		}
		return 0, err
	}

	return gauge, nil
}

func (strg *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	query := `select delta from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(ctx, query, name, types.Counter)
	if row.Err() != nil {
		return 0, fmt.Errorf(exceptions.UnknownMetricName + " " + row.Err().Error())
	}

	var counter int64
	err := row.Scan(&counter)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, fmt.Errorf(exceptions.UnknownMetricName)
		}
		return 0, err
	}

	return counter, nil
}

func (strg *DBStorage) GetGaugeOrZero(ctx context.Context, name string) float64 {
	query := `select value from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(ctx, query, name, types.Gauge)
	if row.Err() != nil {
		return 0
	}

	var gauge float64
	err := row.Scan(&gauge)
	if err != nil {
		return 0
	}

	return gauge
}

func (strg *DBStorage) GetCounterOrZero(ctx context.Context, name string) int64 {
	query := `select delta from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(ctx, query, name, types.Counter)
	if row.Err() != nil {
		return 0
	}

	var counter int64
	err := row.Scan(&counter)
	if err != nil {
		return 0
	}

	return counter
}

func (strg *DBStorage) GetAllGauges(ctx context.Context) map[string]float64 {
	query := `select name, value from metrics where type = $2`
	gauges := make(map[string]float64, 0)

	rows, err := strg.DBConnection.QueryContext(ctx, query, types.Gauge)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return gauges
		}
		return nil
	}

	if rows.Err() != nil {
		return nil
	}

	var gaugeName string
	var gaugeValue float64
	for rows.Next() {
		err := rows.Scan(&gaugeName, &gaugeValue)
		if err != nil {
			return nil
		}
		gauges[gaugeName] = gaugeValue
	}

	return gauges
}

func (strg *DBStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	query := `select name, delta from metrics where type = $2`
	counters := make(map[string]int64, 0)

	rows, err := strg.DBConnection.QueryContext(ctx, query, types.Counter)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return counters
		}
		return nil
	}

	if rows.Err() != nil {
		return nil
	}

	var counterName string
	var counterValue int64
	for rows.Next() {
		err := rows.Scan(&counterName, &counterValue)
		if err != nil {
			return nil
		}
		counters[counterName] = counterValue
	}

	return counters
}

func (strg *DBStorage) CloseDB() error {
	err := strg.DBConnection.Close()
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not close db", err.Error())
	}
	return err
}

func (strg *DBStorage) GetDB() *DBStorage {
	return strg
}

func (strg *DBStorage) GetMemStorage() *MemStorage {
	return nil
}
