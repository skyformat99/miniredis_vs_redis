package main

import (
	"reflect"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/daaku/go.redis/redistest"
	"github.com/garyburd/redigo/redis"
)

type command struct {
	cmd   string // 'GET', 'SET', &c.
	args  []interface{}
	error bool // Whether the command should return an error or not.
}

func succ(cmd string, args ...interface{}) command {
	return command{
		cmd:   cmd,
		args:  args,
		error: false,
	}
}

func fail(cmd string, args ...interface{}) command {
	return command{
		cmd:   cmd,
		args:  args,
		error: true,
	}
}

func TestSet(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		fail("SET", "foo"),
		fail("SET", "foo", "bar", "baz"),
		succ("GET", "foo"),
		succ("SET", "foo", "bar\bbaz"),
		succ("GET", "foo"),
		succ("SET", "foo", "bar", "EX", 100),
		fail("SET", "foo", "bar", "EX", "noint"),
		succ("SET", "utf8", "❆❅❄☃"),
	)
}

func TestGetrange(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "The quick brown fox jumps over the lazy dog"),
		succ("GETRANGE", "foo", 0, 100),
		succ("GETRANGE", "foo", 0, 0),
		succ("GETRANGE", "foo", 0, -4),
		succ("GETRANGE", "foo", 0, -400),
		succ("GETRANGE", "foo", -4, -4),
		succ("GETRANGE", "foo", 4, 2),
		fail("GETRANGE", "foo", "aap", 2),
		fail("GETRANGE", "foo", 4, "aap"),
		fail("GETRANGE", "foo", 4, 2, "aap"),
		fail("GETRANGE", "foo"),
		succ("HSET", "aap", "noot", "mies"),
		fail("GETRANGE", "aap", 4, 2),
	)
}

func testCommands(t *testing.T, commands ...command) {
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()

	sReal, _ := redistest.NewServerClient(t)
	defer sReal.Close()

	cMini, err := redis.Dial("tcp", sMini.Addr())
	ok(t, err)

	cReal, err := redis.Dial(sReal.Proto(), sReal.Addr())
	ok(t, err)

	for _, p := range commands {
		vReal, errReal := cReal.Do(p.cmd, p.args...)
		vMini, errMini := cMini.Do(p.cmd, p.args...)
		if p.error {
			if errReal == nil {
				t.Errorf("got no error from realredis. case: %#v\n", p)
			}
			if errMini == nil {
				t.Errorf("got no error from miniredis. case: %#v\n", p)
			}
		} else {
			if errReal != nil {
				t.Errorf("got an error from realredis: %v. case: %#v\n", errReal, p)
			}
			if errMini != nil {
				t.Errorf("got an error from miniredis: %v. case: %#v\n", errMini, p)
			}
		}
		if !reflect.DeepEqual(errReal, errMini) {
			t.Errorf("error error. expected: %#v got: %#v case: %#v\n", vReal, vMini, p)
		}
		if !reflect.DeepEqual(vReal, vMini) {
			t.Errorf("value error. expected: %#v got: %#v case: %#v\n", vReal, vMini, p)
		}
	}
}

/*
func TestGet(t *testing.T) {
	sMini, err := miniredis.Run()
	ok(t, err)
	defer sMini.Close()

	sReal, _ := redistest.NewServerClient(t)
	defer sReal.Close()

	cMini, err := net.Dial("tcp", sMini.Addr())
	ok(t, err)
	cReal, err := net.Dial("tcp", sReal.Addr())
	ok(t, err)

	c := []string{"SET", "foo", "bar"}
	cMini.Write(bulk(c...))
	bufMini := make([]byte, 1000)
	_, err = cMini.Read(bufMini)
	ok(t, err)
	cReal.Write(bulk(c...))
	bufReal := make([]byte, 1000)
	_, err = cReal.Read(bufReal)
	ok(t, err)
	equals(t, bufReal, bufMini)
}

// Commands to redis 'bulk string' format.
func bulk(cs ...string) []byte {
	res := fmt.Sprintf("*%d\r\n", len(cs))
	for _, c := range cs {
		res += fmt.Sprintf("$%d\r\n%s\r\n", len(c), c)
	}
	return []byte(res)
}
*/
