package errorHandler

import (
	"net/http"
	"runtime"
	auth "chat/handlers/middleware"
)

type ErrorHandler struct {
	function   string
	Err        error
	statusCode int
}

func (e *ErrorHandler) setFuncName() {
	pc, _, _, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()
	e.function = function
}

func (e *ErrorHandler) ResponseForError(w http.ResponseWriter, codeStatus int) {
	e.setFuncName()
	w.WriteHeader(codeStatus)
	http.Error(w, e.Err.Error(), codeStatus)
}

func HandleAuthError(w http.ResponseWriter, r *http.Request, funcName string) {
	auth.DeleteCookies(w)
	http.Redirect(w, r, "/login/", 302)
}
