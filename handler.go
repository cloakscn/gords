package main

import (
	"fmt"
	"sync"
)

// 定义常量
const (
	PING    = "PING"    // ping命令
	SET     = "SET"     // 设置键值对命令
	GET     = "GET"     // 获取键值对命令
	HSET    = "HSET"    // 设置哈希表键值对命令
	HGET    = "HGET"    // 获取哈希表键值对命令
	HGETALL = "HGETALL" // 获取哈希表所有键值对命令
)

// 定义Handlers变量，用于存储命令和对应的处理函数
var Handlers = map[string]func([]Value) Value{
	PING:    ping,    // ping命令对应的处理函数
	SET:     set,     // 设置键值对命令对应的处理函数
	GET:     get,     // 获取键值对命令对应的处理函数
	HSET:    hset,    // 设置哈希表键值对命令对应的处理函数
	HGET:    hget,    // 获取哈希表键值对命令对应的处理函数
	HGETALL: hgetall, // 获取哈希表所有键值对命令对应的处理函数
}

// ping is a command that replies with the argument, or "PONG" if no argument is given.
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: STRING.Str, Str: "PONG"}
	}

	return Value{Typ: STRING.Str, Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

// set sets the value associated with the specified key to the provided value.
// It returns the string reply "OK".
//
// Parameters:
//   - args: A slice of Value, expected to contain exactly two elements:
//     1. The key to set.
//     2. The value to set for the key.
//
// Returns:
//   - Value: A Value struct containing the result of the operation:
//   - If the key exists, returns a simple string reply "OK".
//   - If the number of arguments is incorrect, returns an error message.
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

// get retrieves the value associated with the specified key from the SETs map.
// If the key exists, it returns the corresponding value. If the key doesn't exist,
// it returns a null reply.
//
// Parameters:
//   - args: A slice of Value, expected to contain exactly one element:
//     The key whose value is to be retrieved.
//
// Returns:
//   - Value: A Value struct containing the result of the operation:
//   - If the key exists, returns a bulk string reply with the value.
//   - If the key doesn't exist, returns a null reply.
//   - If the number of arguments is incorrect, returns an error message.
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

// hset sets the field in the hash stored at the specified key to the provided value.
// If the key does not exist, a new key holding a hash is created.
// If the field already exists in the hash, it is overwritten.
//
// Parameters:
//   - args: A slice of Value, expected to contain exactly three elements:
//     1. The key under which the hash is stored.
//     2. The field within the hash to set.
//     3. The value to set for the field.
//
// Returns:
//   - Value: A Value struct containing the result of the operation:
//   - If successful, returns a simple string reply "OK".
//   - If the number of arguments is incorrect, returns an error message.
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

// hget retrieves the value associated with a field in a hash stored at the specified key.
// It returns the value of the field if it exists, or a null reply if either the key or the field do not exist.
//
// Parameters:
//   - args: A slice of Value, expected to contain exactly two elements:
//     1. The key under which the hash is stored.
//     2. The field within the hash whose value is to be retrieved.
//
// Returns:
//   - Value: A Value struct containing the result of the operation:
//   - If successful, returns a bulk string reply with the value of the field.
//   - If the key or field doesn't exist, returns a null reply.
//   - If the number of arguments is incorrect, returns an error message.
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

// hgetall retrieves all key-value pairs from a hash stored at the specified key.
// It returns an array of strings, where each pair of consecutive elements
// represents a key and its value from the hash.
//
// Parameters:
//   - args: A slice of Value, expected to contain exactly one element
//     representing the key of the hash to retrieve.
//
// Returns:
//   - Value: A Value struct containing the result of the operation.
//     If the hash exists, it returns an array of strings with all key-value pairs.
//     If the hash doesn't exist, it returns a null value.
//     If the number of arguments is incorrect, it returns an error message.
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
