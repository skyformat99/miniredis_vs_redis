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
