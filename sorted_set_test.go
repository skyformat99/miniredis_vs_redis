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
		succ("ZREVRANK", "z", "aap"),
		succ("ZREVRANK", "z", "noot"),
		succ("ZREVRANK", "z", "mies"),
		succ("ZREVRANK", "z", "vuur"),
		succ("ZREVRANK", "z", "nosuch"),
		succ("ZREVRANK", "nosuch", "nosuch"),

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
		fail("ZREVRANK"),
		fail("ZREVRANK", "key"),

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
		succ("ZREVRANGE", "z", 0, -1),
		succ("ZREVRANGE", "z", 0, -1, "WITHSCORES"),
		succ("ZREVRANGE", "z", 0, -1, "WiThScOrEs"),
		succ("ZREVRANGE", "z", 0, -2),
		succ("ZREVRANGE", "z", 0, -1000),
		succ("ZREVRANGE", "z", 2, -2),
		succ("ZREVRANGE", "z", 400, -1),
		succ("ZREVRANGE", "z", 300, -110),

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

		fail("ZREVRANGE"),
		fail("ZREVRANGE", "str", 300, -110),
	)
}

func TestSortedSetRem(t *testing.T) {
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
		succ("ZREM", "z", "nosuch"),
		succ("ZREM", "z", "mies", "nootagain"),
		succ("ZRANGE", "z", 0, -1),

		// failure cases
		fail("ZREM"),
		fail("ZREM", "foo"),
		succ("SET", "str", "I am a string"),
		fail("ZREM", "str", "member"),
	)
}

func TestSortedSetScore(t *testing.T) {
	testCommands(t,
		succ("ZADD", "z",
			1, "aap",
			2, "noot",
			3, "mies",
			2, "nootagain",
			3, "miesagain",
			math.Inf(+1), "the stars",
		),
		succ("ZSCORE", "z", "mies"),
		succ("ZSCORE", "z", "the stars"),
		succ("ZSCORE", "z", "nosuch"),
		succ("ZSCORE", "nosuch", "nosuch"),

		// failure cases
		fail("ZSCORE"),
		fail("ZSCORE", "foo"),
		fail("ZSCORE", "foo", "too", "many"),
		succ("SET", "str", "I am a string"),
		fail("ZSCORE", "str", "member"),
	)
}

func TestSortedSetRangeByScore(t *testing.T) {
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
		succ("ZRANGEBYSCORE", "z", "-inf", "inf"),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 1, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", -1, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 1, -2),
		succ("ZREVRANGEBYSCORE", "z", "inf", "-inf"),
		succ("ZREVRANGEBYSCORE", "z", "inf", "-inf", "LIMIT", 1, 2),
		succ("ZREVRANGEBYSCORE", "z", "inf", "-inf", "LIMIT", -1, 2),
		succ("ZREVRANGEBYSCORE", "z", "inf", "-inf", "LIMIT", 1, -2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "WITHSCORES"),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "WiThScOrEs"),
		succ("ZREVRANGEBYSCORE", "z", "-inf", "inf", "WITHSCORES", "LIMIT", 1, 2),
		succ("ZRANGEBYSCORE", "z", 0, 3),
		succ("ZRANGEBYSCORE", "z", 0, "inf"),
		succ("ZRANGEBYSCORE", "z", "(1", "3"),
		succ("ZRANGEBYSCORE", "z", "(1", "(3"),
		succ("ZRANGEBYSCORE", "z", "1", "(3"),
		succ("ZRANGEBYSCORE", "z", "1", "(3", "LIMIT", 0, 2),
		succ("ZRANGEBYSCORE", "foo", 2, 3, "LIMIT", 1, 2, "WITHSCORES"),
		succ("ZCOUNT", "z", "-inf", "inf"),
		succ("ZCOUNT", "z", 0, 3),
		succ("ZCOUNT", "z", 0, "inf"),
		succ("ZCOUNT", "z", "(2", "inf"),

		// Bunch of limit edge cases
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 0, 7),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 0, 8),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 0, 9),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 7, 0),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 7, 1),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 7, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 8, 0),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 8, 1),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 8, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", 9, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", -1, 2),
		succ("ZRANGEBYSCORE", "z", "-inf", "inf", "LIMIT", -1, -1),

		// failure cases
		fail("ZRANGEBYSCORE"),
		fail("ZRANGEBYSCORE", "foo"),
		fail("ZRANGEBYSCORE", "foo", 1),
		fail("ZRANGEBYSCORE", "foo", 2, 3, "toomany"),
		fail("ZRANGEBYSCORE", "foo", 2, 3, "WITHSCORES", "toomany"),
		fail("ZRANGEBYSCORE", "foo", 2, 3, "LIMIT", "noint", 1),
		fail("ZRANGEBYSCORE", "foo", 2, 3, "LIMIT", 1, "noint"),
		fail("ZREVRANGEBYSCORE", "z", "-inf", "inf", "WITHSCORES", "LIMIT", 1, -2, "toomany"),
		fail("ZRANGEBYSCORE", "foo", "noint", 3),
		fail("ZRANGEBYSCORE", "foo", "[4", 3),
		fail("ZRANGEBYSCORE", "foo", 2, "noint"),
		fail("ZRANGEBYSCORE", "foo", "4", "[3"),
		succ("SET", "str", "I am a string"),
		fail("ZRANGEBYSCORE", "str", 300, -110),

		fail("ZREVRANGEBYSCORE"),
		fail("ZREVRANGEBYSCORE", "foo", "[4", 3),
		fail("ZREVRANGEBYSCORE", "str", 300, -110),

		fail("ZCOUNT"),
		fail("ZCOUNT", "foo", "[4", 3),
		fail("ZCOUNT", "str", 300, -110),
	)
}
