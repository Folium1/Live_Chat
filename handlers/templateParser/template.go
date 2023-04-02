package templateparser

import (
	"chat/logger"
	"html/template"
)

func TemplateParse() *template.Template {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't parse tamplate", err, funcName)
	}
	return tmpl
}
