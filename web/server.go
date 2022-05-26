package web

import (
	"database/sql"
	"net/http"
	"text/template"
)

var Tpl = template.Must(template.ParseGlob("templates/*.html"))

type myServer struct {
	Db     *sql.DB
	Router *http.ServeMux
	serve  *http.Server
}
