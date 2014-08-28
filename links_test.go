package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetHashFromId(t *testing.T) {
	Convey("A id should be able to hash and unhash transparently", t, func() {
		hash := GetHashFromId(1042)
		So(hash, ShouldEqual, "sy")
		id := GetIdFromHash("sy")
		So(id, ShouldEqual, 1042)
	})
}

func TestRedisLinkerCreateLink(t *testing.T) {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", "127.0.0.1:6379")
	}, 10)

	redisLinker := &RedisLinker{redisPool}

	conn := redisPool.Get()
	conn.Do("FLUSHALL")
	conn.Close()

	Convey("A link should always be able to created, and return a proper ID", t, func() {
		link, err := redisLinker.CreateLink(Link{
			Default: "http://www.google.com/",
		})
		So(err, ShouldBeNil)
		So(link, ShouldNotBeNil)

		So(link.Id, ShouldEqual, 1)
		So(link.Default, ShouldEqual, "http://www.google.com/")
		So(link.Android, ShouldEqual, "")
		So(link.IOS, ShouldEqual, "")

		conn := redisPool.Get()
		linkId, err := redis.Int(conn.Do("GET", "link_id"))
		So(err, ShouldBeNil)
		So(linkId, ShouldEqual, 1)
	})
}
