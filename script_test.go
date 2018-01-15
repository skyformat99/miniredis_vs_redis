package main

// Script commands

import (
	"testing"
)

func TestEval(t *testing.T) {
	testCommands(t,
		succ("EVAL", "return 42", 0),
		succ("EVAL", "", 0),
		succ("EVAL", "return 42", 1, "foo"),
		succ("EVAL", "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", 2, "key1", "key2", "first", "second"),
		succ("EVAL", "return {ARGV[1]}", 0, "first"),
		succ("EVAL", "return {ARGV[1]}", 0, "first\nwith\nnewlines!\r\r\n\t!"),

		// failure cases
		fail("EVAL"),
		fail("EVAL", "return 42"),
		fail("EVAL", "["),
		fail("EVAL", "return 42", "return 43"),
		fail("EVAL", "return 42", 1),
		fail("EVAL", "return 42", -1),
		fail("EVAL", 42),
	)
}

func TestScript(t *testing.T) {
	testCommands(t,
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return 43"),

		succ("SCRIPT", "EXISTS", "1fa00e76656cc152ad327c13fe365858fd7be306"),
		succ("SCRIPT", "EXISTS", "0", "1fa00e76656cc152ad327c13fe365858fd7be306"),
		succ("SCRIPT", "EXISTS", 0),
		succ("SCRIPT", "EXISTS"),

		succ("SCRIPT", "FLUSH"),
		succ("SCRIPT", "EXISTS", "1fa00e76656cc152ad327c13fe365858fd7be306"),

		fail("SCRIPT"),
		fail("SCRIPT", "LOAD", "return 42", "return 42"),
		failLoosely("SCRIPT", "LOAD", "]"),
		fail("SCRIPT", "LOAD", "]", "foo"),
		fail("SCRIPT", "LOAD"),
		fail("SCRIPT", "FLUSH", "foo"),
		fail("SCRIPT", "FOO"),
	)
}

func TestEvalsha(t *testing.T) {
	sha1 := "1fa00e76656cc152ad327c13fe365858fd7be306" // "return 42"
	sha2 := "bfbf458525d6a0b19200bfd6db3af481156b367b" // keys[1], argv[1]

	testCommands(t,
		succ("SCRIPT", "LOAD", "return 42"),
		succ("SCRIPT", "LOAD", "return {KEYS[1],ARGV[1]}"),
		succ("EVALSHA", sha1, "0"),
		succ("EVALSHA", sha2, "0"),
		succ("EVALSHA", sha2, "0", "foo"),
		succ("EVALSHA", sha2, "1", "foo"),
		succ("EVALSHA", sha2, "1", "foo", "bar"),
		succ("EVALSHA", sha2, "1", "foo", "bar", "baz"),

		succ("SCRIPT", "FLUSH"),
		fail("EVALSHA", sha1, "0"),

		succ("SCRIPT", "LOAD", "return 42"),
		fail("EVALSHA", sha1),
		fail("EVALSHA"),
		fail("EVALSHA", "nosuch"),
		fail("EVALSHA", "nosuch", 0),
	)
}
