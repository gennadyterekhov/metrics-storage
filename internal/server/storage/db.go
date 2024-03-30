package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"reflect"
)

type DBStorage struct {
	DBConnection       *sql.DB
	HTTPRequestContext context.Context
}

func CreateDBStorage() *DBStorage {
	conn, err := sql.Open("pgx", config.Conf.DBDsn)
	if err != nil {
		logger.ZapSugarLogger.Panicln("could not connect to db using dsn: " + config.Conf.DBDsn)
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

func (strg *DBStorage) SetContext(ctx context.Context) {
	strg.HTTPRequestContext = ctx
}

func (strg *DBStorage) getContext() context.Context {
	if strg.HTTPRequestContext == nil {
		logger.ZapSugarLogger.Debugln("DB does not have context")
		return context.Background()
	}
	return strg.HTTPRequestContext
}

func (strg *DBStorage) Clear() {
	_, err := strg.DBConnection.Exec("delete from metrics")
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
	}
}

func (strg *DBStorage) AddCounter(key string, value int64) {
	query := `
			INSERT INTO metrics ( name, type, delta)
			VALUES ($1, $2, $3)
			ON CONFLICT(name) 
			DO UPDATE SET
			  delta = metrics.delta + $3;`

	_, err := strg.DBConnection.ExecContext(strg.getContext(), query, key, types.Counter, value)
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not add counter", err.Error())
	}
}

func (strg *DBStorage) SetGauge(key string, value float64) {
	query := `
			INSERT INTO metrics ( name, type, value)
			VALUES ($1, $2, $3)
			ON CONFLICT(name) 
			DO UPDATE SET
			  value = $3;`

	_, err := strg.DBConnection.ExecContext(strg.getContext(), query, key, types.Gauge, value)
	if err != nil {
		logger.ZapSugarLogger.Errorln(err.Error())
	}
}

func (strg *DBStorage) GetGauge(name string) (float64, error) {

	query := `select value from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(strg.getContext(), query, name, types.Gauge)
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

func (strg *DBStorage) GetCounter(name string) (int64, error) {
	query := `select delta from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(strg.getContext(), query, name, types.Counter)
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

func (strg *DBStorage) GetGaugeOrZero(name string) float64 {
	query := `select value from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(strg.getContext(), query, name, types.Gauge)
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

func (strg *DBStorage) GetCounterOrZero(name string) int64 {
	query := `select delta from metrics where name = $1 and type = $2`

	row := strg.DBConnection.QueryRowContext(strg.getContext(), query, name, types.Counter)
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

func (strg *DBStorage) GetAllGauges() map[string]float64 {
	query := `select name, value from metrics where type = $2`
	var gauges = make(map[string]float64, 0)

	rows, err := strg.DBConnection.QueryContext(strg.getContext(), query, types.Gauge)
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

func (strg *DBStorage) GetAllCounters() map[string]int64 {
	query := `select name, delta from metrics where type = $2`
	var counters = make(map[string]int64, 0)

	rows, err := strg.DBConnection.QueryContext(strg.getContext(), query, types.Counter)
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

func (strg *DBStorage) IsEqual(anotherStorage StorageInterface) (eq bool) {
	gauges, counters := strg.GetAllGauges(), strg.GetAllCounters()
	gauges2, counters2 := anotherStorage.GetAllGauges(), anotherStorage.GetAllCounters()

	return reflect.DeepEqual(gauges, gauges2) && reflect.DeepEqual(counters, counters2)
}

func (strg *DBStorage) CloseDB() error {
	err := strg.DBConnection.Close()
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not close db", err.Error())
	}
	return err
}

func (strg *DBStorage) IsDB() bool {
	return true
}

func (strg *DBStorage) GetDB() *DBStorage {
	return strg
}
