package main

// Sorted Set keys.

import (
	"math"
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

		succ("RENAME", "z", "z2"),
		succ("EXISTS", "z"),
		succ("EXISTS", "z2"),
		succ("MOVE", "z2", 3),
		succ("EXISTS", "z2"),
		succ("SELECT", 3),
		succ("EXISTS", "z2"),
		succ("DEL", "z2"),
		succ("EXISTS", "z2"),
	)
}

func TestSortedSetRange(t *testing.T) {
	testCommands(t,
		succ("ZADD", "z",
			1, "aap",
			2, "noot",
			3, "mies",
			2, "nootagain",
			3, "miesagain",
			math.Inf(+1), "the stars",
			math.Inf(+1), "more stars",
			math.Inf(-1), "big bang",
		),
		succ("ZRANGE", "z", 0, -1),
		succ("ZRANGE", "z", 0, -1, "WITHSCORES"),
		succ("ZRANGE", "z", 0, -1, "WiThScOrEs"),
		succ("ZRANGE", "z", 0, -2),
		succ("ZRANGE", "z", 0, -1000),
		succ("ZRANGE", "z", 2, -2),
		succ("ZRANGE", "z", 400, -1),
		succ("ZRANGE", "z", 300, -110),

		succ("ZADD", "zz",
			0, "aap",
			0, "Aap",
			0, "AAP",
			0, "aAP",
			0, "aAp",
		),
		succ("ZRANGE", "zz", 0, -1),

		// failure cases
		fail("ZRANGE"),
		fail("ZRANGE", "foo"),
		fail("ZRANGE", "foo", 1),
		fail("ZRANGE", "foo", 2, 3, "toomany"),
		fail("ZRANGE", "foo", 2, 3, "WITHSCORES", "toomany"),
		fail("ZRANGE", "foo", "noint", 3),
		fail("ZRANGE", "foo", 2, "noint"),
		succ("SET", "str", "I am a string"),
		fail("ZRANGE", "str", 300, -110),
	)
}
