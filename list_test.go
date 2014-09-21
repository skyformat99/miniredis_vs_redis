package main

// List keys.

import (
	"testing"
)

func TestLPush(t *testing.T) {
	testCommands(t,
		succ("LPUSH", "l", "aap", "noot", "mies"),
		succ("TYPE", "l"),
		succ("LPUSH", "l", "more", "keys"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "l", 0, 6),
		succ("LRANGE", "l", 2, 6),
		succ("LRANGE", "l", -100, -100),
		succ("LRANGE", "nosuch", 2, 6),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("LPOP", "l"),
		succ("EXISTS", "l"),
		succ("LPOP", "nosuch"),

		// failure cases
		fail("LPUSH"),
		fail("LPUSH", "l"),
		succ("SET", "str", "I am a string"),
		fail("LPUSH", "str", "noot", "mies"),
		fail("LRANGE"),
		fail("LRANGE", "key"),
		fail("LRANGE", "key", 2),
		fail("LRANGE", "key", 2, 6, "toomany"),
		fail("LRANGE", "key", "noint", 6),
		fail("LRANGE", "key", 2, "noint"),
		fail("LPOP"),
		fail("LPOP", "key", "args"),
	)
}

func TestRPush(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("TYPE", "l"),
		succ("RPUSH", "l", "more", "keys"),
		succ("LRANGE", "l", 0, -1),
		succ("LRANGE", "l", 0, 6),
		succ("LRANGE", "l", 2, 6),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("RPOP", "l"),
		succ("EXISTS", "l"),
		succ("RPOP", "nosuch"),

		// failure cases
		fail("RPUSH"),
		fail("RPUSH", "l"),
		succ("SET", "str", "I am a string"),
		fail("RPUSH", "str", "noot", "mies"),
		fail("RPOP"),
		fail("RPOP", "key", "args"),
	)
}

func TestLinxed(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LINDEX", "l", 0),
		succ("LINDEX", "l", 1),
		succ("LINDEX", "l", 2),
		succ("LINDEX", "l", 3),
		succ("LINDEX", "l", 4),
		succ("LINDEX", "l", 44444),
		succ("LINDEX", "l", -0),
		succ("LINDEX", "l", -1),
		succ("LINDEX", "l", -2),
		succ("LINDEX", "l", -3),
		succ("LINDEX", "l", -4),
		succ("LINDEX", "l", -4000),

		// failure cases
		fail("LINDEX"),
		fail("LINDEX", "l"),
		succ("SET", "str", "I am a string"),
		fail("LINDEX", "str", 1),
		fail("LINDEX", "l", "noint"),
		fail("LINDEX", "l", 1, "too many"),
	)
}

func TestLlen(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LLEN", "l"),
		succ("LLEN", "nosuch"),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("LLEN", "str"),
		fail("LLEN"),
		fail("LLEN", "l", "too many"),
	)
}

func TestLtrim(t *testing.T) {
	testCommands(t,
		succ("RPUSH", "l", "aap", "noot", "mies"),
		succ("LTRIM", "l", 0, 1),
		succ("LRANGE", "l", 0, -1),
		succ("RPUSH", "l2", "aap", "noot", "mies", "vuur"),
		succ("LTRIM", "l2", -2, -1),
		succ("LRANGE", "l2", 0, -1),
		succ("RPUSH", "l3", "aap", "noot", "mies", "vuur"),
		succ("LTRIM", "l3", -2, -1000),
		succ("LRANGE", "l3", 0, -1),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("LTRIM", "str", 0, 1),
		fail("LTRIM", "l", 0, 1, "toomany"),
		fail("LTRIM", "l", "noint", 1),
		fail("LTRIM", "l", 0, "noint"),
		fail("LTRIM", "l", 0),
		fail("LTRIM", "l"),
		fail("LTRIM"),
	)
}
