package main

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type Link struct {
	Id      int    `redis:"id"`
	Default string `redis:"default"`
	Android string `redis:"android"`
	IOS     string `redis:"ios"`
}

func GetHashFromId(id int) string {
	return strconv.FormatInt(int64(id), 36)
}

func GetIdFromHash(hash string) int {
	id, _ := strconv.ParseInt(hash, 36, 0)
	return int(id)
}

type Linker interface {
	CreateLink(Link) (*Link, error)
	FindLinkById(id int) (*Link, error)
}

type RedisLinker struct {
	pool *redis.Pool
}

func (rl *RedisLinker) CreateLink(link Link) (*Link, error) {
	conn := rl.pool.Get()
	defer conn.Close()

	linkId, err := redis.Int(conn.Do("INCR", "link_id"))
	if err != nil {
		return nil, err
	}

	conn.Do("HMSET", fmt.Sprintf("link:%d", linkId),
		"id", linkId,
		"default", link.Default,
		"android", link.Android,
		"ios", link.IOS,
	)

	return rl.FindLinkById(linkId)
}

func (rl *RedisLinker) FindLinkById(id int) (*Link, error) {
	conn := rl.pool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("link:%d", id)))
	if err != nil {
		return nil, err
	}

	if len(reply) <= 0 {
		return nil, errors.New("Invalid link hash requested!")
	}

	link := &Link{}

	if err := redis.ScanStruct(reply, link); err != nil {
		return nil, err
	}

	return link, nil
}
