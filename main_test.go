package main

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/daaku/go.redis/redistest"
	"github.com/garyburd/redigo/redis"
)

func TestSet(t *testing.T) {
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()

	sReal, _ := redistest.NewServerClient(t)
	defer sReal.Close()

	cMini, err := redis.Dial("tcp", sMini.Addr())
	ok(t, err)

	cReal, err := redis.Dial(sReal.Proto(), sReal.Addr())
	ok(t, err)

	{
		vMini, errMini := cMini.Do("SET", "foo", "bar")
		vReal, errReal := cReal.Do("SET", "foo", "bar")
		equals(t, errReal, errMini)
		equals(t, vReal, vMini)
	}

	{
		vMini, errMini := cMini.Do("SET", "foo")
		vReal, errReal := cReal.Do("SET", "foo")
		equals(t, errReal, errMini)
		equals(t, vReal, vMini)
	}
}
