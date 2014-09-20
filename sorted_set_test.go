package main

// Sorted Set keys.

import (
	"testing"
)

func TestSortedSet(t *testing.T) {
	testCommands(t,
		succ("ZADD", "z", 1, "aap", 2, "noot", 3, "mies"),
		succ("ZADD", "z", 1, "vuur", 4, "noot"),
		succ("TYPE", "z"),
		succ("EXISTS", "z"),
		succ("ZCARD", "z"),

		succ("ZRANK", "z", "aap"),
		succ("ZRANK", "z", "noot"),
		succ("ZRANK", "z", "mies"),
		succ("ZRANK", "z", "vuur"),
		succ("ZRANK", "z", "nosuch"),
		succ("ZRANK", "nosuch", "nosuch"),

		succ("ZADD", "zi", "inf", "aap", "-inf", "noot", "+inf", "mies"),
		succ("ZRANK", "zi", "noot"),

		// Double key
		succ("ZADD", "zz", 1, "aap", 2, "aap"),
		succ("ZCARD", "zz"),

		// failure cases
		succ("SET", "str", "I am a string"),
		fail("ZADD"),
		fail("ZADD", "s"),
		fail("ZADD", "s", 1),
		fail("ZADD", "s", 1, "aap", 1),
		fail("ZADD", "s", "nofloat", "aap"),
		fail("ZADD", "str", 1, "aap"),
		fail("ZCARD"),
		fail("ZCARD", "too", "many"),
		fail("ZCARD", "str"),
		fail("ZRANK"),
		fail("ZRANK", "key"),
		fail("ZRANK", "key", "too", "many"),
		fail("ZRANK", "str", "member"),
	)
}
