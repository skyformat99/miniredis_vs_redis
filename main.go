// Package to compare the Redis implementation of github.com/alicebob/miniredis
// against a real Redis server.
//
// There is no executable code beside the tests.
//
package main

import (
	_ "github.com/alicebob/miniredis"
	_ "github.com/garyburd/redigo/redis"
)

// We're only interested in the tests.

func main() {
}
