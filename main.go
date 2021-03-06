/*

Source code of the Moje Zdrowie

*/

package main

import (
	"html/template"
	"moje-zdrowie/downloader"
	"moje-zdrowie/medicines"
	"moje-zdrowie/medmap"
	"moje-zdrowie/pharmacies"
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
		port = "80"
	}

	go downloader.StartDownloads()

	http.HandleFunc("/", start)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/map", medmap.Page)
	http.HandleFunc("/medicines", medicines.Page)
	http.HandleFunc("/refunded", refunded.Page)
	http.HandleFunc("/pharmacy", pharmacies.GetPharmacyInfo)

	http.ListenAndServe(":"+port, nil)
}
