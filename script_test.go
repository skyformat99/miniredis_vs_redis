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

func TestLua(t *testing.T) {
	// basic datatype things
	testCommands(t,
		succ("EVAL", "", 0),
		succ("EVAL", "return 42", 0),
		succ("EVAL", "return 42, 43", 0),
		succ("EVAL", "return true", 0),
		succ("EVAL", "return 'foo'", 0),
		succ("EVAL", "return 3.1415", 0),
		succ("EVAL", "return 3.9999", 0),
		succ("EVAL", "return {1,'foo'}", 0),
		succ("EVAL", "return {1,'foo',nil,'foo'}", 0),
		succ("EVAL", "return 3.9999+3", 0),
		succ("EVAL", "return 3.99+0.0001", 0),
		succ("EVAL", "return 3.9999+0.201", 0),
		succ("EVAL", "return {{1}}", 0),
		succ("EVAL", "return {1,{1,{1,'bar'}}}", 0),
	)
	// special returns
	testCommands(t,
		fail("EVAL", "return {err = 'oops'}", 0),
		succ("EVAL", "return {1,{err = 'oops'}}", 0),
		fail("EVAL", "return redis.error_reply('oops')", 0),
		succ("EVAL", "return {1,redis.error_reply('oops')}", 0),
		fail("EVAL", "return {err = 'oops', noerr = true}", 0), // doc error?
		fail("EVAL", "return {1, 2, err = 'oops'}", 0),         // doc error?

		succ("EVAL", "return {ok = 'great'}", 0),
		succ("EVAL", "return {1,{ok = 'great'}}", 0),
		succ("EVAL", "return redis.status_reply('great')", 0),
		succ("EVAL", "return {1,redis.status_reply('great')}", 0),
		succ("EVAL", "return {ok = 'great', notok = 'yes'}", 0),       // doc error?
		succ("EVAL", "return {1, 2, ok = 'great', notok = 'yes'}", 0), // doc error?

		failLoosely("EVAL", "return redis.error_reply(1)", 0),
		failLoosely("EVAL", "return redis.error_reply()", 0),
		failLoosely("EVAL", "return redis.error_reply(redis.error_reply('foo'))", 0),
		failLoosely("EVAL", "return redis.status_reply(1)", 0),
		failLoosely("EVAL", "return redis.status_reply()", 0),
		failLoosely("EVAL", "return redis.status_reply(redis.status_reply('foo'))", 0),
	)
	// state inside lua
	testCommands(t,
		succ("EVAL", "redis.call('SELECT', 3); redis.call('SET', 'foo', 'bar')", 0),
		succ("GET", "foo"),
		succ("SELECT", 3),
		succ("GET", "foo"),
	)
}
