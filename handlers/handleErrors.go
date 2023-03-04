package handlers

import (
	"log"
	"net/http"
	"runtime"
)

type errorHandler struct {
	function   string
	err        error
	statusCode int
}

func (e *errorHandler) setFuncName() {
	pc, _, _, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()
	e.function = function
}

func (e *errorHandler) responseForError(w http.ResponseWriter, codeStatus int) {
	e.setFuncName()
	log.Printf("Function:%v | error: %v | code: %v", e.function, e.err, e.statusCode)
	w.WriteHeader(codeStatus)
	http.Error(w, e.err.Error(), codeStatus)
}
