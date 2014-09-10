package main

// Transaction things.

import (
	"testing"
)

func TestTx(t *testing.T) {
	testCommands(t,
		succ("MULTI"),
		succ("SET", "AAP", 1),
		succ("GET", "AAP"),
		succ("EXEC"),
		succ("GET", "AAP"),
	)

	// err: Double MULTI
	testCommands(t,
		succ("MULTI"),
		fail("MULTI"),
	)

	// err: No MULTI
	testCommands(t,
		fail("EXEC"),
	)

	// Errors in the MULTI sequence
	testCommands(t,
		succ("MULTI"),
		succ("SET", "foo", "bar"),
		fail("SET", "foo"),
		succ("SET", "foo", "bar"),
		fail("EXEC"),
	)

	// Simple WATCH
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("WATCH", "foo"),
		succ("MULTI"),
		succ("GET", "foo"),
		succ("EXEC"),
	)

	// Simple UNWATCH
	testCommands(t,
		succ("SET", "foo", "bar"),
		succ("WATCH", "foo"),
		succ("UNWATCH"),
		succ("MULTI"),
		succ("GET", "foo"),
		succ("EXEC"),
	)

	// UNWATCH in a MULTI. Yep. Weird.
	testCommands(t,
		succ("WATCH", "foo"),
		succ("MULTI"),
		succ("UNWATCH"), // Valid. Somehow.
		succ("EXEC"),
	)
}
