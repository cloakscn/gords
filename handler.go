package main

import (
	"fmt"
	"sync"

	"github.com/cloakscn/gords/message"
)

var Handlers = map[string]func([]message.Value) message.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []message.Value) message.Value {
	if len(args) == 0 {
		return message.Value{Typ: message.STRING.Str, Str: "PONG"}
	}

	return message.Value{Typ: message.STRING.Str, Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []message.Value) message.Value {
	if len(args) != 2 {
		return message.Value{Typ: message.ERROR.Str, Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return message.Value{Typ: message.STRING.Str, Str: "OK"}
}

func get(args []message.Value) message.Value {
	if len(args) != 1 {
		return message.Value{Typ: message.ERROR.Str, Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return message.Value{Typ: message.NULL.Str}
	}

	return message.Value{Typ: message.BULK.Str, Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []message.Value) message.Value {
	if len(args) != 3 {
		return message.Value{Typ: message.ERROR.Str, Str: "ERR wrong number of arguments for 'hset' command"}
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

	return message.Value{Typ: message.STRING.Str, Str: "OK"}
}

// hget retrieves the value associated with the specified key from the hash map
// identified by the given hash. It expects exactly two arguments: the hash and
// the key. If the key does not exist within the hash, a null value is returned.
// If the number of arguments is incorrect, it returns an error.
func hget(args []message.Value) message.Value {
	if len(args) != 2 {
		return message.Value{Typ: message.ERROR.Str, Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return message.Value{Typ: message.NULL.Str}
	}

	return message.Value{Typ: message.BULK.Str, Bulk: value}
}

func hgetall(args []message.Value) message.Value {
	if len(args) != 1 {
		return message.Value{Typ: message.BULK.Str, Str: "ERR wrong number of arguments for 'hgetall' commandmessage."}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	values, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return message.Value{Typ: message.NULL.Str}
	}

	result := message.Value{Typ: message.ARRAY.Str}
	for key, value := range values {
		result.Array = append(result.Array, message.Value{
			Typ: message.BULK.Str, Bulk: fmt.Sprintf("%s: %s", key, value)})
	}

	return result
}
