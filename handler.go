package main

import (
	"fmt"
	"sync"

)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: STRING.Str, Str: "PONG"}
	}

	return Value{Typ: STRING.Str, Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: ERROR.Str, Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{Typ: STRING.Str, Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: ERROR.Str, Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{Typ: NULL.Str}
	}

	return Value{Typ: BULK.Str, Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Typ: ERROR.Str, Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{Typ: STRING.Str, Str: "OK"}
}

// hget retrieves the value associated with the specified key from the hash map
// identified by the given hash. It expects exactly two arguments: the hash and
// the key. If the key does not exist within the hash, a null value is returned.
// If the number of arguments is incorrect, it returns an error.
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: ERROR.Str, Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: NULL.Str}
	}

	return Value{Typ: BULK.Str, Bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: BULK.Str, Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	values, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: NULL.Str}
	}

	result := Value{Typ: ARRAY.Str}
	for key, value := range values {
		result.Array = append(result.Array, Value{
			Typ: BULK.Str, Bulk: fmt.Sprintf("%s: %s", key, value)})
	}

	return result
}
