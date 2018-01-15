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
