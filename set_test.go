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
		succSorted("SMEMBERS", "s"),
		succSorted("SMEMBERS", "nosuch"),

		// failure cases
		fail("SADD"),
		fail("SADD", "s"),
		fail("SMEMBERS"),
		// Wrong type
		succ("SET", "str", "I am a string"),
		fail("SADD", "str", "noot", "mies"),
		fail("SMEMBERS", "str"),

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
