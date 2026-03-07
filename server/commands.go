package server

import (
	"github.com/Pujan-khunt/redis-lite/resp"
	"github.com/Pujan-khunt/redis-lite/storage"
)

// CommandHandler defines the signature for any function which handles a Redis command.
type CommandHandler = func(args []resp.RespValue, w *resp.RespWriter, store storage.Store)

var commandRegistry = map[string]CommandHandler{
	"SET":  handleSET,
	"GET":  handleGET,
	"DEL":  handleDEL,
	"PING": handlePING,
}

func handleSET(args []resp.RespValue, w *resp.RespWriter, store storage.Store) {
	if len(args) != 3 {
		w.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'SET' command"})
		return
	}
	key, val := args[1].Str, args[2].Str
	store.Set(key, val)
	w.Write(resp.RespValue{Type: resp.SimpleString, Str: "OK"})
}

func handleGET(args []resp.RespValue, w *resp.RespWriter, store storage.Store) {
	if len(args) != 2 {
		w.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'GET' command"})
		return
	}
	key := args[1].Str
	if val, ok := store.Get(key); ok {
		w.Write(resp.RespValue{Type: resp.BulkString, Str: val})
	} else {
		w.Write(resp.RespValue{Type: resp.BulkString, Str: ""})
	}
}

func handleDEL(args []resp.RespValue, w *resp.RespWriter, store storage.Store) {
	if len(args) != 2 {
		w.Write(resp.RespValue{Type: resp.Error, Str: "-ERR invalid number of arguments for 'DEL' command\r\n"})
		return
	}
	key := args[1].Str
	if ok := store.Del(key); ok {
		w.Write(resp.RespValue{Type: resp.Integer, Num: 1})
	} else {
		w.Write(resp.RespValue{Type: resp.Error, Str: "-ERR failed to delete key"})
	}
}

func handlePING(args []resp.RespValue, w *resp.RespWriter, store storage.Store) {
	if len(args) == 2 && args[1].Type == resp.BulkString {
		// PONG [arg] for PING [arg]
		w.Write(resp.RespValue{Type: resp.BulkString, Str: args[1].Str})
	} else {
		// PONG for PING
		w.Write(resp.RespValue{Type: resp.SimpleString, Str: "PONG"})
	}
}
