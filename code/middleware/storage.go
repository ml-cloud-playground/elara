package main

import "fmt"

// Storage is a wrapper for combined cache and database operations
type Storage struct {
	sqlstorage SQLStorage
	cache      *Cache
}

// Init kicks off the database connector
func (s *Storage) Init(user, password, host, name, redisHost, redisPort string, cache bool) error {
	if err := s.sqlstorage.Init(user, password, host, name); err != nil {
		return err
	}

	var err error
	s.cache, err = NewCache(redisHost, redisPort, cache)
	if err != nil {
		return err
	}

	return nil
}

// Match returns a single nonprofit from cache or database
func (s Storage) Match(subcategory string) (nonprofit, error) {
	t, err := s.sqlstorage.Read(subcategory)
	if err != nil {
		return t, fmt.Errorf("error getting single from database elara: %v", err)
	}
	return t, nil
}
