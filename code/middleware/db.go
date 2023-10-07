package main

import (
	"database/sql"
)

// SQLStorage is a wrapper for database operations
type SQLStorage struct {
	db *sql.DB
}

// Init kicks off the database connector
func (s *SQLStorage) Init(user, password, host, name string) error {
	var err error
	s.db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+")/"+name+"?parseTime=true")
	if err != nil {
		return err
	}

	return nil
}

// Close ends the database connection
func (s *SQLStorage) Close() error {
	return s.db.Close()
}

// Match returns matching non-profit based on subcategory
func (s SQLStorage) Read(subcategory string) (nonprofit, error) {
	t := nonprofit{}
	results, err := s.db.Query("SELECT * FROM nonprofits WHERE subcategory =? order by RAND() LIMIT 1", subcategory)
	if err != nil {
		return t, err
	}

	results.Next()
	t, err = resultToNonProfit(results)
	if err != nil {
		return t, err
	}

	return t, nil
}

func resultToNonProfit(results *sql.Rows) (nonprofit, error) {
	t := nonprofit{}
	if err := results.Scan(&t.ID, &t.Name, &t.SubCategory); err != nil {
		return t, err
	}
	return t, nil
}
