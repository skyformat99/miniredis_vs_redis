package main

// Set keys.

import (
	"testing"
)

func TestSet(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SADD", "s", "vuur", "noot"),
		succ("TYPE", "s"),
		succ("EXISTS", "s"),
		succ("SCARD", "s"),
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "nosuch"),
		succ("SISMEMBER", "s", "aap"),
		succ("SISMEMBER", "s", "nosuch"),

		succ("SCARD", "nosuch"),
		succ("SISMEMBER", "nosuch", "nosuch"),

		// failure cases
		fail("SADD"),
		fail("SADD", "s"),
		fail("SMEMBERS"),
		fail("SMEMBERS", "too", "many"),
		fail("SCARD"),
		fail("SCARD", "too", "many"),
		fail("SISMEMBER"),
		fail("SISMEMBER", "few"),
		fail("SISMEMBER", "too", "many", "arguments"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SADD", "str", "noot", "mies"),
		fail("SMEMBERS", "str"),
		fail("SISMEMBER", "str", "noot"),
		fail("SCARD", "str"),
	)
}

func TestSetMove(t *testing.T) {
	// Move a set around
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("RENAME", "s", "others"),
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "others"),
		succ("MOVE", "others", 2),
		succSorted("SMEMBERS", "others"),
		succ("SELECT", 2),
		succSorted("SMEMBERS", "others"),
	)
}

func TestSetDel(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SREM", "s", "noot", "nosuch"),
		succ("SCARD", "s"),
		succSorted("SMEMBERS", "s"),

		// failure cases
		fail("SREM"),
		fail("SREM", "s"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SREM", "str", "noot"),
	)
}

func TestSetSMove(t *testing.T) {
	testCommands(t,
		succ("SADD", "s", "aap", "noot", "mies"),
		succ("SMOVE", "s", "s2", "aap"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),
		succ("SMOVE", "s", "s2", "nosuch"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),
		succ("SMOVE", "s", "nosuch", "noot"),
		succ("SCARD", "s"),
		succ("SCARD", "s2"),

		succ("SMOVE", "s", "s2", "mies"),
		succ("SCARD", "s"),
		succ("EXISTS", "s"),
		succ("SCARD", "s2"),
		succ("EXISTS", "s2"),

		succ("SMOVE", "s2", "s2", "mies"),

		succ("SADD", "s5", "aap"),
		succ("SADD", "s6", "aap"),
		succ("SMOVE", "s5", "s6", "aap"),

		// failure cases
		fail("SMOVE"),
		fail("SMOVE", "s"),
		fail("SMOVE", "s", "s2"),
		fail("SMOVE", "s", "s2", "too", "many"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SMOVE", "str", "s2", "noot"),
		fail("SMOVE", "s2", "str", "noot"),
	)
}

func TestSetSpop(t *testing.T) {
	testCommands(t,
		// Set with a single member...
		succ("SADD", "s", "aap"),
		succ("SPOP", "s"),
		succ("EXISTS", "s"),

		succ("SPOP", "nosuch"),

		// failure cases
		fail("SPOP"),
		fail("SPOP", "s", "s2"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SPOP", "str"),
	)
}

func TestSetSrandmember(t *testing.T) {
	testCommands(t,
		// Set with a single member...
		succ("SADD", "s", "aap"),
		succ("SRANDMEMBER", "s"),
		succ("SRANDMEMBER", "s", 1),
		succ("SRANDMEMBER", "s", 5),
		succ("SRANDMEMBER", "s", -1),
		succ("SRANDMEMBER", "s", -5),

		succ("SRANDMEMBER", "s", 0),
		succ("SPOP", "nosuch"),

		// failure cases
		fail("SRANDMEMBER"),
		fail("SRANDMEMBER", "s", "noint"),
		fail("SRANDMEMBER", "s", 1, "toomany"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SRANDMEMBER", "str"),
	)
}
