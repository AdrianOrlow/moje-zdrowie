package medmap

import (
	"html/template"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	t, _ := template.ParseFiles("templates/medmap/medmap.html")
	t.Execute(w, nil)
}
