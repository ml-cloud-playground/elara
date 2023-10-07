package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// RedisPool is an interface that allows us to swap in an mock for testing cache
// code.
type RedisPool interface {
	Get() redis.Conn
}

// ErrCacheMiss error indicates that an item is not in the cache
var ErrCacheMiss = fmt.Errorf("item is not in cache")

// NewCache returns an initialized cache ready to go.
func NewCache(redisHost, redisPort string, enabled bool) (*Cache, error) {
	c := &Cache{}
	pool := c.InitPool(redisHost, redisPort)
	c.enabled = enabled
	c.redisPool = pool
	return c, nil
}

// Cache abstracts all of the operations of caching for the application
type Cache struct {
	// redisPool *redis.Pool
	redisPool RedisPool
	enabled   bool
}

func (c *Cache) log(msg string) {
	log.Printf("Cache     : %s\n", msg)
}

// InitPool starts the cache off
func (c Cache) InitPool(redisHost, redisPort string) RedisPool {
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	msg := fmt.Sprintf("Initialized Redis at %s", redisAddr)
	c.log(msg)
	const maxConnections = 10

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, maxConnections)

	return pool
}

// Clear removes all items from the cache.
func (c Cache) Clear() error {
	if !c.enabled {
		return nil
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	if _, err := conn.Do("FLUSHALL"); err != nil {
		return err
	}
	return nil
}

// Save records a nonprofit into the cache.
func (c *Cache) Save(nonprofit nonprofit) error {
	if !c.enabled {
		return nil
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	json, err := nonprofit.JSON()
	if err != nil {
		return fmt.Errorf("cannot convert nonprofit to json: %s", err)
	}

	conn.Send("MULTI")
	conn.Send("SET", strconv.Itoa(nonprofit.ID), json)

	if _, err := conn.Do("EXEC"); err != nil {
		return fmt.Errorf("cannot perform exec operation on cache: %s", err)
	}
	c.log("Successfully saved nonprofit to cache")
	return nil
}

// Get gets a nonprofit from the cache.
func (c *Cache) Get(key string) (nonprofit, error) {
	t := nonprofit{}
	if !c.enabled {
		return t, ErrCacheMiss
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nonprofit{}, ErrCacheMiss
	} else if err != nil {
		return nonprofit{}, err
	}

	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return nonprofit{}, err
	}
	c.log("Successfully retrieved nonprofit from cache")

	return t, nil
}

// Delete will remove a nonprofit from the cache completely.
func (c *Cache) Delete(key string) error {
	if !c.enabled {
		return nil
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}

	c.log(fmt.Sprintf("Cleaning from cache %s", key))
	return nil
}

// List gets all of the nonprofits from the cache.
func (c *Cache) List() (nonprofits, error) {
	t := nonprofits{}
	if !c.enabled {
		return t, ErrCacheMiss
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", "nonprofitslist"))
	if err == redis.ErrNil {
		return nonprofits{}, ErrCacheMiss
	} else if err != nil {
		return nonprofits{}, err
	}

	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return nonprofits{}, err
	}
	c.log("Successfully retrieved nonprofits from cache")

	return t, nil
}

// SaveList records a nonprofit list into the cache.
func (c *Cache) SaveList(nonprofits nonprofits) error {
	if !c.enabled {
		return nil
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	json, err := nonprofits.JSON()
	if err != nil {
		return err
	}

	if _, err := conn.Do("SET", "nonprofitslist", json); err != nil {
		return err
	}
	c.log("Successfully saved nonprofit to cache")
	return nil
}

// DeleteList deletes a nonprofit list into the cache.
func (c *Cache) DeleteList() error {
	if !c.enabled {
		return nil
	}

	return c.Delete("nonprofitslist")
}
