package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Row struct {
	*sql.Row
}

//TODO: support other types
//Raw.scan wrapper allows scan *bool, *string, *int64 types
func (row Row) ScanNill(dest ...interface{}) error {
	replaceDest := make([]interface{}, 0)
	for _, val := range dest {
		switch reflect.TypeOf(val) {
		case reflect.PtrTo(reflect.TypeOf(true)):
			replaceDest = append(replaceDest, new(sql.NullBool))
		case reflect.PtrTo(reflect.TypeOf("")):
			replaceDest = append(replaceDest, new(sql.NullString))
		case reflect.PtrTo(reflect.TypeOf(int64(1))):
			replaceDest = append(replaceDest, new(sql.NullInt64))
		default:
			replaceDest = append(replaceDest, val)
		}
	}
	if err := row.Scan(replaceDest...); err != nil {
		return err
	}

	for i, val := range replaceDest {
		//fmt.Println(reflect.TypeOf(val))
		switch reflect.TypeOf(val) {
		case reflect.TypeOf(new(sql.NullInt64)):
			t := val.(*sql.NullInt64)
			if t.Valid {
				*(dest[i].(*int64)) = t.Int64
			}
		case reflect.TypeOf(new(sql.NullBool)):
			t := val.(*sql.NullBool)
			if t.Valid {
				*(dest[i].(*bool)) = t.Bool
			}
		case reflect.TypeOf(new(sql.NullString)):
			t := val.(*sql.NullString)
			if t.Valid {
				*(dest[i].(*string)) = t.String
			}
		}
	}

	return nil
}

//DBLogger is wrapper for sql.DB that log Exec and Query calls
type DBLogger struct {
	*sql.DB
	Log *log.Logger
}

//genId generate random id. Id is used for identify which query invoke the error
func genId() int {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	return rand.Int()
}

func replaceRule(args []interface{}) func(string) string {
	l := len(args)
	return func(in string) string {
		n, _ := strconv.Atoi(in[1:])
		if n > l {
			return "<Error arg>"
		}
		return fmt.Sprintf("'%v'", args[n-1])
	}
}

func formatQuery(query string, args ...interface{}) string {
	reg, _ := regexp.Compile(`\$\d+`)
	return reg.ReplaceAllStringFunc(query, replaceRule(args))
}

func (dbl *DBLogger) Exec(query string, args ...interface{}) (*ResWrp, error) {

	dbl.Log.WithFields(log.Fields{"Log_id": genId(), "query": formatQuery(query, args...)}).Infof("client request")

	res, err := dbl.DB.Exec(query, args...)
	if err != nil {
		log.Warnln(err)
	}
	return &ResWrp{res}, err
}

//Wrapper for sql.Result provide some auxiliary functions for checking correctness of result of query
type ResWrp struct {
	sql.Result
}

var ErrMuchRows = errors.New("Affected rows more than expected")

//Check that affected row is just one
func (res *ResWrp) AffectedOnlyRow() error {
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return sql.ErrNoRows
	}
	if i > 1 {
		return ErrMuchRows
	}
	return nil
}

func (res *ResWrp) AffectedAtLeastRow() error {
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (dbl *DBLogger) Query(query string, args ...interface{}) (*RowsLogger, error) {
	dbl.Log.WithFields(log.Fields{"Log_id": genId(), "query": formatQuery(query, args...)}).Infof("client request")

	rows, err := dbl.DB.Query(query, args...)
	if err != nil {
		log.Warnln(err)
	}
	return &RowsLogger{Rows: rows}, err
}

type RowsLogger struct {
	Rows *sql.Rows
	Log  *log.Entry
}

func (rsl *RowsLogger) Scan(dest ...interface{}) error {
	err := rsl.Rows.Scan(dest...)
	if err != nil {
		rsl.Log.Warnln(err)
	}
	return err
}

func (dbl *DBLogger) QueryRow(query string, args ...interface{}) *RowLogger {
	dbl.Log.WithFields(log.Fields{"Log_id": genId(), "query": formatQuery(query, args...)}).Infof("client request")

	return &RowLogger{Row: dbl.DB.QueryRow(query, args...)}
}

type RowLogger struct {
	Row *sql.Row
	log *log.Entry
}

func (rl *RowLogger) Scan(dest ...interface{}) error {
	err := rl.Row.Scan(dest...)
	if err != nil {
		rl.log.Warnln(err)
	}
	return err
}
