package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func StartDownloads() {
	for {
		sources := map[string]string{
			"downloads/pharmacies.xml": "http://pub.rejestrymedyczne.csioz.gov.pl/Pobieranie_WS/Pobieranie.ashx?filetype=XMLFile&regtype=RA_FILES",
			"downloads/medicines.xml":  "http://pub.rejestrymedyczne.csioz.gov.pl/pobieranie_WS/Pobieranie.ashx?filetype=XMLFile&regtype=RPL_FILES_BASE",
		}
		/* pharmacies.ConvertPharmaciesToMarkers() */
		for filepath, url := range sources {
			go DownloadFile(filepath, url)
		}
		time.Sleep(6 * time.Hour)
	}
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
