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
