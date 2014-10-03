package main

import (
	"testing"
)

func TestKeys(t *testing.T) {
	testCommands(t,
		succ("SET", "one", "1"),
		succ("SET", "two", "2"),
		succ("SET", "three", "3"),
		succ("SET", "four", "4"),
		succSorted("KEYS", `*o*`),
		succSorted("KEYS", `t??`),
		succSorted("KEYS", `t?*`),
		succSorted("KEYS", `*`),
		succSorted("KEYS", `t*`),
		succSorted("KEYS", `t\*`),
		succSorted("KEYS", `[tf]*`),

		// zero length key
		succ("SET", "", "nothing"),
		succ("GET", ""),

		// Simple failure cases
		fail("KEYS"),
		fail("KEYS", "foo", "bar"),
	)

	testCommands(t,
		succ("SET", "[one]", "1"),
		succ("SET", "two", "2"),
		succSorted("KEYS", `[\[o]*`),
		succSorted("KEYS", `\[*`),
		succSorted("KEYS", `*o*`),
		succSorted("KEYS", `[]*`), // nothing
	)
}

func TestRandom(t *testing.T) {
	testCommands(t,
		succ("RANDOMKEY"),
		// A random key from a DB with a single key. We can test that.
		succ("SET", "one", "1"),
		succ("RANDOMKEY"),

		// Simple failure cases
		fail("RANDOMKEY", "bar"),
	)
}

func TestUnknownCommand(t *testing.T) {
	// Can't compare; we get a different message from redeo
	testCommands(t,
		fail("nosuch"), // redeo doesn't change the capitilization, Redis lowercases it.
		succ("SET", "foo", "bar"),
	)
}

func TestQuit(t *testing.T) {
	testCommands(t,
		succ("QUIT"),
		fail("QUIT"),
	)
}

func TestRename(t *testing.T) {
	testCommands(t,
		// No 'a' key
		fail("RENAME", "a", "b"),

		// Move a key with the TTL.
		succ("SET", "a", "3"),
		succ("EXPIRE", "a", "123"),
		succ("SET", "b", "12"),
		succ("RENAME", "a", "b"),
		succ("EXISTS", "a"),
		succ("GET", "a"),
		succ("TYPE", "a"),
		succ("TTL", "a"),
		succ("EXISTS", "b"),
		succ("GET", "b"),
		succ("TYPE", "b"),
		succ("TTL", "b"),

		// Error cases
		fail("RENAME"),
		fail("RENAME", "a"),
		fail("RENAME", "a", "b", "toomany"),
	)
}
