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

func TestSet(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("GET", "foo"),
		succ("SET", "foo", "bar\bbaz"),
		succ("GET", "foo"),
		succ("SET", "foo", "bar", "EX", 100),
		fail("SET", "foo", "bar", "EX", "noint"),
		succ("SET", "utf8", "❆❅❄☃"),

		succ("GETSET", "foo", "new"),
		succ("GET", "foo"),
		// GETSET on a new key
		succ("GETSET", "nosuch", "new"),
		succ("GET", "nosuch"),
		// Failure cases
		fail("SET"),
		fail("SET", "foo"),
		fail("SET", "foo", "bar", "baz"),
		succ("HSET", "hash", "key", "value"),
		fail("GET", "hash"),
		fail("GETSET", "hash", "new"),
	)
}

func TestExpire(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("EXPIRE", "foo", 12),
		succ("TTL", "foo"),
		succ("TTL", "nosuch"),
		succ("SET", "foo", "bar"),
		succ("PEXPIRE", "foo", 999999),
		succ("EXPIREAT", "foo", 2234567890),
		succ("PEXPIREAT", "foo", 2234567890000),
		// succ("PTTL", "foo"),
		succ("PTTL", "nosuch"),

		fail("EXPIRE"),
		fail("EXPIRE", "foo"),
		fail("EXPIRE", "foo", "noint"),
		fail("EXPIRE", "foo", 12, "toomany"),
		fail("EXPIREAT"),
		fail("TTL"),
		fail("TTL", "too", "many"),
		fail("PEXPIRE"),
		fail("PEXPIRE", "foo"),
		fail("PEXPIRE", "foo", "noint"),
		fail("PEXPIRE", "foo", 12, "toomany"),
		fail("PEXPIREAT"),
		fail("PTTL"),
		fail("PTTL", "too", "many"),
	)
}

func TestMset(t *testing.T) {
	testCommands(t,
		succ("MSET", "foo", "bar"),
		succ("MSET", "foo", "bar", "baz", "?"),
		succ("MSET", "foo", "bar", "foo", "baz"), // double key
		succ("GET", "foo"),
		// Error cases
		fail("MSET"),
		fail("MSET", "foo"),
		fail("MSET", "foo", "bar", "baz"),

		succ("MSETNX", "foo", "bar", "aap", "noot"),
		succ("MSETNX", "one", "two", "three", "four"),
		succ("MSETNX", "11", "12", "11", "14"), // double key
		succ("GET", "11"),
		// Error cases
		fail("MSETNX"),
		fail("MSETNX", "one"),
		fail("MSETNX", "one", "two", "three"),

		succ("HSET", "aap", "noot", "mies"),
		succ("MSETNX", "aap", "again", "eight", "nine"),
	)
}

func TestSetx(t *testing.T) {
	testCommands(t,
		succ("SETEX", "foo", 12, "bar"),
		succ("GET", "foo"),
		succ("TTL", "foo"),
		fail("SETEX", "foo"),
		fail("SETEX", "foo", "noint", "bar"),
		fail("SETEX", "foo", 12),
		fail("SETEX", "foo", 12, "bar", "toomany"),

		succ("PSETEX", "foo", 12, "bar"),
		succ("GET", "foo"),
		// succ("PTTL", "foo"), // counts down too quickly to compare
		fail("PSETEX", "foo"),
		fail("PSETEX", "foo", "noint", "bar"),
		fail("PSETEX", "foo", 12),
		fail("PSETEX", "foo", 12, "bar", "toomany"),
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

func TestStrlen(t *testing.T) {
	testCommands(t,
		succ("SET", "str", "The quick brown fox jumps over the lazy dog"),
		succ("STRLEN", "str"),
		// failure cases
		fail("STRLEN"),
		fail("STRLEN", "str", "bar"),
		succ("HSET", "hash", "key", "value"),
		fail("STRLEN", "hash"),
	)
}

func TestSetrange(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "The quick brown fox jumps over the lazy dog"),
		succ("SETRANGE", "foo", 0, "aap"),
		succ("GET", "foo"),
		succ("SETRANGE", "foo", 10, "noot"),
		succ("GET", "foo"),
		succ("SETRANGE", "foo", 40, "overtheedge"),
		succ("GET", "foo"),
		succ("SETRANGE", "foo", 400, "oh, hey there"),
		succ("GET", "foo"),
		// Non existing key
		succ("SETRANGE", "nosuch", 2, "aap"),
		succ("GET", "nosuch"),

		// Error cases
		fail("SETRANGE", "foo"),
		fail("SETRANGE", "foo", 1),
		fail("SETRANGE", "foo", "aap", "bar"),
		fail("SETRANGE", "foo", "noint", "bar"),
		fail("SETRANGE", "foo", -1, "bar"),
		succ("HSET", "aap", "noot", "mies"),
		fail("SETRANGE", "aap", 4, "bar"),
	)
}

func TestIncrAndFriends(t *testing.T) {
	testCommands(t,
		succ("INCR", "aap"),
		succ("INCR", "aap"),
		succ("INCR", "aap"),
		succ("GET", "aap"),
		succ("DECR", "aap"),
		succ("DECR", "noot"),
		succ("DECR", "noot"),
		succ("GET", "noot"),
		succ("INCRBY", "noot", 100),
		succ("INCRBY", "noot", 200),
		succ("INCRBY", "noot", 300),
		succ("GET", "noot"),
		succ("DECRBY", "noot", 100),
		succ("DECRBY", "noot", 200),
		succ("DECRBY", "noot", 300),
		succ("DECRBY", "noot", 400),
		succ("GET", "noot"),
		succ("INCRBYFLOAT", "zus", 1.23),
		succ("INCRBYFLOAT", "zus", 3.1456),
		succ("INCRBYFLOAT", "zus", 987.65432),
		succ("GET", "zus"),
		succ("INCRBYFLOAT", "whole", 300),
		succ("INCRBYFLOAT", "whole", 300),
		succ("INCRBYFLOAT", "whole", 300),
		succ("GET", "whole"),
		succ("INCRBYFLOAT", "big", 12345e10),
		succ("GET", "big"),

		// Floats are not ints.
		succ("SET", "float", 1.23),
		fail("INCR", "float"),
		fail("INCRBY", "float", 12),
		fail("DECR", "float"),
		fail("DECRBY", "float", 12),
		succ("SET", "str", "I'm a string"),
		fail("INCRBYFLOAT", "str", 123.5),

		// Error cases
		succ("HSET", "mies", "noot", "mies"),
		fail("INCR", "mies"),
		fail("INCRBY", "mies", 1),
		fail("INCRBY", "mies", "foo"),
		fail("DECR", "mies"),
		fail("DECRBY", "mies", 1),
		fail("INCRBYFLOAT", "mies", 1),
		fail("INCRBYFLOAT", "int", "foo"),

		fail("INCR", "int", "err"),
		fail("INCRBY", "int"),
		fail("DECR", "int", "err"),
		fail("DECRBY", "int"),
		fail("INCRBYFLOAT", "int"),
	)
}

func TestBitcount(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "The quick brown fox jumps over the lazy dog"),
		succ("SET", "utf8", "❆❅❄☃"),
		succ("BITCOUNT", "foo"),
		succ("BITCOUNT", "utf8"),
		succ("BITCOUNT", "foo", 0, 0),
		succ("BITCOUNT", "utf8", 0, 0),
		fail("BITCOUNT", "foo", 4, 2, 2, 2, 2),
		succ("HSET", "aap", "noot", "mies"),
		fail("BITCOUNT", "aap", 4, 2),
	)
}

func TestBitop(t *testing.T) {
	testCommands(t,
		succ("SET", "a", "foo"),
		succ("SET", "b", "aap"),
		succ("SET", "c", "noot"),
		succ("SET", "d", "mies"),
		succ("SET", "e", "❆❅❄☃"),

		// ANDs
		succ("BITOP", "AND", "target", "a", "b", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "AND", "target", "a", "nosuch", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "AND", "utf8", "e", "e"),
		succ("GET", "utf8"),
		succ("BITOP", "AND", "utf8", "b", "e"),
		succ("GET", "utf8"),
		// BITOP on only unknown keys:
		succ("BITOP", "AND", "bits", "nosuch", "nosucheither"),
		succ("GET", "bits"),

		// ORs
		succ("BITOP", "OR", "target", "a", "b", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "OR", "target", "a", "nosuch", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "OR", "utf8", "e", "e"),
		succ("GET", "utf8"),
		succ("BITOP", "OR", "utf8", "b", "e"),
		succ("GET", "utf8"),
		// BITOP on only unknown keys:
		succ("BITOP", "OR", "bits", "nosuch", "nosucheither"),
		succ("GET", "bits"),
		succ("SET", "empty", ""),
		// BITOP on empty key
		succ("BITOP", "OR", "bits", "empty"),
		succ("GET", "bits"),

		// XORs
		succ("BITOP", "XOR", "target", "a", "b", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "XOR", "target", "a", "nosuch", "c", "d"),
		succ("GET", "target"),
		succ("BITOP", "XOR", "target", "a"),
		succ("GET", "target"),
		succ("BITOP", "XOR", "utf8", "e", "e"),
		succ("GET", "utf8"),
		succ("BITOP", "XOR", "utf8", "b", "e"),
		succ("GET", "utf8"),

		// NOTs
		succ("BITOP", "NOT", "target", "a"),
		succ("GET", "target"),
		succ("BITOP", "NOT", "target", "e"),
		succ("GET", "target"),
		succ("BITOP", "NOT", "bits", "nosuch"),
		succ("GET", "bits"),

		fail("BITOP", "AND", "utf8"),
		fail("BITOP", "AND"),
		fail("BITOP", "NOT", "foo", "bar", "baz"),
		fail("BITOP", "WRONGOP", "key"),
		fail("BITOP", "WRONGOP"),

		succ("HSET", "hash", "aap", "noot"),
		fail("BITOP", "AND", "t", "hash", "irrelevant"),
		fail("BITOP", "OR", "t", "hash", "irrelevant"),
		fail("BITOP", "XOR", "t", "hash", "irrelevant"),
		fail("BITOP", "NOT", "t", "hash"),
	)
}

func TestBitpos(t *testing.T) {
	testCommands(t,
		succ("SET", "a", "\x00\x0f"),
		succ("SET", "b", "\xf0\xf0"),
		succ("SET", "c", "\x00\x00\x00\x0f"),
		succ("SET", "d", "\x00\x00\x00"),
		succ("SET", "e", "\xff\xff\xff"),

		succ("BITPOS", "a", 1),
		succ("BITPOS", "a", 0),
		succ("BITPOS", "a", 1, 1),
		succ("BITPOS", "a", 0, 1),
		succ("BITPOS", "a", 1, 1, 2),
		succ("BITPOS", "a", 0, 1, 2),
		succ("BITPOS", "b", 1),
		succ("BITPOS", "b", 0),
		succ("BITPOS", "c", 1),
		succ("BITPOS", "c", 0),
		succ("BITPOS", "d", 1),
		succ("BITPOS", "d", 0),
		succ("BITPOS", "e", 1),
		succ("BITPOS", "e", 0),
		succ("BITPOS", "e", 1, 1),
		succ("BITPOS", "e", 0, 1),
		succ("BITPOS", "e", 1, 1, 2),
		succ("BITPOS", "e", 0, 1, 2),
		succ("BITPOS", "e", 1, 100, 2),
		succ("BITPOS", "e", 0, 100, 2),
		succ("BITPOS", "e", 1, 1, -2),
		succ("BITPOS", "e", 1, 1, -2000),
		succ("BITPOS", "e", 0, 1, 2),
		succ("BITPOS", "nosuch", 1),
		succ("BITPOS", "nosuch", 0),

		succ("HSET", "hash", "aap", "noot"),
		fail("BITPOS", "hash", 1),
		fail("BITPOS", "a", "aap"),
	)
}

func TestGetbit(t *testing.T) {
	commands := []command{
		succ("SET", "a", "\x00\x0f"),
		succ("SET", "e", "\xff\xff\xff"),
		succ("GETBIT", "nosuch", 1),
		succ("GETBIT", "nosuch", 0),

		// Error cases
		succ("HSET", "hash", "aap", "noot"),
		fail("GETBIT", "hash", 1),
		fail("GETBIT", "a", "aap"),
		fail("GETBIT", "a"),
		fail("GETBIT", "too", 1, "many"),
	}

	// Generate read commands.
	for i := range make([]struct{}, 100) {
		commands = append(commands,
			succ("GETBIT", "a", i),
			succ("GETBIT", "e", i),
		)
	}

	testCommands(t, commands...)
}

func TestSetbit(t *testing.T) {
	commands := []command{
		succ("SET", "a", "\x00\x0f"),
		succ("SETBIT", "a", 0, 1),
		succ("GET", "a"),
		succ("SETBIT", "a", 0, 0),
		succ("GET", "a"),
		succ("SETBIT", "a", 13, 0),
		succ("GET", "a"),
		succ("SETBIT", "nosuch", 11111, 1),
		succ("GET", "nosuch"),

		// Error cases
		succ("HSET", "hash", "aap", "noot"),
		fail("SETBIT", "hash", 1, 1),
		fail("SETBIT", "a", "aap", 0),
		fail("SETBIT", "a", 0, "aap"),
		fail("SETBIT", "a", -1, 0),
		fail("SETBIT", "a", 1, -1),
		fail("SETBIT", "a", 1, 2),
		fail("SETBIT", "too", 1, 2, "many"),
	}

	// Generate read commands.
	for i := range make([]struct{}, 100) {
		commands = append(commands,
			succ("GETBIT", "a", i),
			succ("GETBIT", "e", i),
		)
	}

	testCommands(t, commands...)
}

func TestAppend(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("APPEND", "foo", "more"),
		succ("GET", "foo"),
		succ("APPEND", "nosuch", "more"),
		succ("GET", "nosuch"),

		// Failure cases
		fail("APPEND"),
		fail("APPEND", "foo"),
	)
}

func TestMove(t *testing.T) {
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("EXPIRE", "foo", 12345),
		succ("MOVE", "foo", 2),
		succ("GET", "foo"),
		succ("TTL", "foo"),
		succ("SELECT", 2),
		succ("GET", "foo"),
		succ("TTL", "foo"),

		// Failure cases
		fail("MOVE"),
		fail("MOVE", "foo"),
		// fail("MOVE", "foo", "noint"),
	)
	// hash key
	testCommands(t,
		succ("HSET", "hash", "key", "value"),
		succ("EXPIRE", "hash", 12345),
		succ("MOVE", "hash", 2),
		succ("MGET", "hash", "key"),
		succ("TTL", "hash"),
		succ("SELECT", 2),
		succ("MGET", "hash", "key"),
		succ("TTL", "hash"),
	)
	testCommands(t,
		succ("SET", "foo", "bar"),
		// to current DB.
		fail("MOVE", "foo", 0),
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
