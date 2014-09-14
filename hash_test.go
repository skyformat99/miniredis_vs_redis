package main

// Hash keys.

import (
	"testing"
)

func TestHash(t *testing.T) {
	testCommands(t,
		succ("HSET", "aap", "noot", "mies"),
		succ("HGET", "aap", "noot"),
		succ("HMGET", "aap", "noot"),
		succ("HLEN", "aap"),
		succ("HKEYS", "aap"),
		succ("HVALS", "aap"),

		succ("HDEL", "aap", "noot"),
		succ("HGET", "aap", "noot"),
		succ("EXISTS", "aap"), // key is gone

		// failure cases
		fail("HSET", "aap", "noot"),
		fail("HGET", "aap"),
		fail("HMGET", "aap"),
		fail("HLEN"),
		fail("HKEYS"),
		fail("HVALS"),
		succ("SET", "str", "I am a string"),
		fail("HSET", "str", "noot", "mies"),
		fail("HGET", "str", "noot"),
		fail("HMGET", "str", "noot"),
		fail("HLEN", "str"),
		fail("HKEYS", "str"),
		fail("HVALS", "str"),
	)
}

func TestHmset(t *testing.T) {
	testCommands(t,
		succ("HMSET", "aap", "noot", "mies", "vuur", "zus"),
		succ("HGET", "aap", "noot"),
		succ("HGET", "aap", "vuur"),
		succ("HLEN", "aap"),

		// failure cases
		fail("HMSET", "aap"),
		fail("HMSET", "aap", "key"),
		fail("HMSET", "aap", "key", "value", "odd"),
		succ("SET", "str", "I am a string"),
		fail("HMSET", "str", "key", "value"),
	)
}

func TestHashIncr(t *testing.T) {
	testCommands(t,
		succ("HINCRBY", "aap", "noot", 12),
		succ("HINCRBY", "aap", "noot", -13),
		succ("HINCRBY", "aap", "noot", 2123),
		succ("HGET", "aap", "noot"),

		// Simple failure cases.
		fail("HINCRBY"),
		fail("HINCRBY", "aap"),
		fail("HINCRBY", "aap", "noot"),
		fail("HINCRBY", "aap", "noot", "noint"),
		fail("HINCRBY", "aap", "noot", 12, "toomany"),
		succ("SET", "str", "value"),
		fail("HINCRBY", "str", "value", 12),
		succ("HINCRBY", "aap", "noot", 12),
	)

	testCommands(t,
		succ("HINCRBYFLOAT", "aap", "noot", 12.3),
		succ("HINCRBYFLOAT", "aap", "noot", -13.1),
		succ("HINCRBYFLOAT", "aap", "noot", 200),
		succ("HGET", "aap", "noot"),

		// Simple failure cases.
		fail("HINCRBYFLOAT"),
		fail("HINCRBYFLOAT", "aap"),
		fail("HINCRBYFLOAT", "aap", "noot"),
		fail("HINCRBYFLOAT", "aap", "noot", "noint"),
		fail("HINCRBYFLOAT", "aap", "noot", 12, "toomany"),
		succ("SET", "str", "value"),
		fail("HINCRBYFLOAT", "str", "value", 12),
		succ("HINCRBYFLOAT", "aap", "noot", 12),
	)
}
