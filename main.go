/*

Source code of the Moje Zdrowie

*/

package main

import (
	"html/template"
	"moje-zdrowie/medicines"
	"moje-zdrowie/medmap"
	"moje-zdrowie/refunded"
	"net/http"
	"os"
)

func start(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/start.html")
	t.Execute(w, nil)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", start)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/map", medmap.Page)
	http.HandleFunc("/medicines", medicines.Page)
	http.HandleFunc("/refunded", refunded.Page)

	http.ListenAndServe(":"+port, nil)
}
