package main

import (
	"testing"
)

func TestEcho(t *testing.T) {
	testCommands(t,
		succ("ECHO", "hello world"),
		fail("ECHO", "hello", "world"),
		fail("eChO", "hello", "world"),
	)
}

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

func TestRenamenx(t *testing.T) {
	testCommands(t,
		// No 'a' key
		fail("RENAMENX", "a", "b"),

		succ("SET", "a", "value"),
		succ("SET", "str", "value"),
		succ("RENAMENX", "a", "str"),
		succ("EXISTS", "a"),
		succ("EXISTS", "str"),
		succ("GET", "a"),
		succ("GET", "str"),

		succ("RENAMENX", "a", "nosuch"),
		succ("EXISTS", "a"),
		succ("EXISTS", "nosuch"),

		// Error cases
		fail("RENAMENX"),
		fail("RENAMENX", "a"),
		fail("RENAMENX", "a", "b", "toomany"),
	)
}

func TestScan(t *testing.T) {
	testCommands(t,
		// No keys yet
		succ("SCAN", 0),

		succ("SET", "key", "value"),
		succ("SCAN", 0),
		succ("SCAN", 0, "COUNT", 12),
		succ("SCAN", 0, "cOuNt", 12),

		succ("SET", "anotherkey", "value"),
		succ("SCAN", 0, "MATCH", "anoth*"),
		succ("SCAN", 0, "MATCH", "anoth*", "COUNT", 100),
		succ("SCAN", 0, "COUNT", 100, "MATCH", "anoth*"),

		// Can't really test multiple keys.
		// succ("SET", "key2", "value2"),
		// succ("SCAN", 0),

		// Error cases
		fail("SCAN"),
		fail("SCAN", "noint"),
		fail("SCAN", 0, "COUNT", "noint"),
		fail("SCAN", 0, "COUNT"),
		fail("SCAN", 0, "MATCH"),
		fail("SCAN", 0, "garbage"),
		fail("SCAN", 0, "COUNT", 12, "MATCH", "foo", "garbage"),
	)
}
