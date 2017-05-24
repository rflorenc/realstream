package main

import "net/http"
import "sync"

// type HandlerFunc func(http.ResponseWriter, *http.Request)

// We will store the map of variables keyed with the
// request instances that the variables belong to.
var vars map[*http.Request]map[string]interface{}
var varsLock sync.RWMutex

func OpenVars(r *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[r] = map[string]interface{}{}
	varsLock.Unlock()
}

// This function safely deletes the entry in the vars map for the request.
func CloseVars(r *http.Request) {
	varsLock.Lock()
	delete(vars, r)
	varsLock.Unlock()
}

// Safely Get and Set our global map properties
func GetVar(r *http.Request, key string) interface{} {
	varsLock.RLock()
	value := vars[r][key]
	varsLock.RUnlock()
	return value
}

func SetVar(r *http.Request, key string, value interface{}) {
	varsLock.Lock()
	vars[r][key] = value
	varsLock.Unlock()
}
