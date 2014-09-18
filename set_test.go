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

		// Move a set around
		succ("RENAME", "s", "others"),
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "others"),
		succ("MOVE", "others", 2),
		succSorted("SMEMBERS", "others"),
		succ("SELECT", 2),
		succSorted("SMEMBERS", "others"),
	)
}
