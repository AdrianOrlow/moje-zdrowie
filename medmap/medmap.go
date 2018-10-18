package medmap

import (
	"html/template"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/medmap/medmap.html")
	t.Execute(w, nil)
}
