package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
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

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
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
				lError(t, "got no error from realredis. case: %#v\n", p)
			}
			if errMini == nil {
				lError(t, "got no error from miniredis. case: %#v\n", p)
			}
		} else {
			if errReal != nil {
				lError(t, "got an error from realredis: %v. case: %#v\n", errReal, p)
			}
			if errMini != nil {
				lError(t, "got an error from miniredis: %v. case: %#v\n", errMini, p)
			}
		}
		if !reflect.DeepEqual(errReal, errMini) {
			lError(t, "error error. expected: %#v got: %#v case: %#v\n",
				vReal, vMini, p)
		}
		if !reflect.DeepEqual(vReal, vMini) {
			lError(t, "value error. expected: %#v got: %#v case: %#v\n",
				vReal, vMini, p)
		}
	}
}

func lError(t *testing.T, format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	prefix := fmt.Sprintf("%s:%d: ", filepath.Base(file), line)
	fmt.Printf(prefix+format, args...)
	t.Fail()
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
