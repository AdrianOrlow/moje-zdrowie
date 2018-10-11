package medicines

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Medicines struct {
	XMLName    xml.Name   `xml:"produktyLecznicze"`
	LastUpdate string     `xml:"stanNaDzien,attr"`
	Medicines  []Medicine `xml:"produktLeczniczy"`
}

type Medicine struct {
	XMLName          xml.Name  `xml:"produktLeczniczy"`
	FullName         string    `xml:"nazwaProduktu,attr"`
	ShortName        string    `xml:"nazwaPowszechnieStosowana,attr"`
	Strength         string    `xml:"moc,attr"`
	Kind             string    `xml:"postac,attr"`
	Producer         string    `xml:"podmiotOdpowiedzialny,attr"`
	Validity         string    `xml:"waznoscPozwolenia,attr"`
	ID               string    `xml:"id,attr"`
	ActiveSubstances []string  `xml:"substancjeCzynne>substancjaCzynna"`
	Packages         []Package `xml:"opakowania>opakowanie"`
}

type Package struct {
	XMLName  xml.Name `xml:"opakowanie"`
	Size     string   `xml:"wielkosc,attr"`
	SizeUnit string   `xml:"jednostkaWielkosci,attr"`
	EanCode  string   `xml:"kodEAN,attr"`
	EuNumber string   `xml:"numerEu,attr"`
}

func MoreInfo(w http.ResponseWriter, r *http.Request, medID string) {
	infoData := Medicine{}

	var medicines Medicines
	medicinesData, _ := os.Open("downloads/medicines.xml")
	defer medicinesData.Close()
	byteValue, _ := ioutil.ReadAll(medicinesData)
	xml.Unmarshal(byteValue, &medicines)

	for i := 0; i < len(medicines.Medicines); i++ {
		for j := 0; j < len(medicines.Medicines[i].Packages); j++ {
			switch strings.Contains(medicines.Medicines[i].ID, medID) {
			case true:
				infoData = medicines.Medicines[i]
				break
			}
		}
	}

	t, _ := template.ParseFiles("templates/medicines/info.html")
	t.Execute(w, infoData)
}

func Results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	if strings.Join(form["medID"], "") != "" {
		MoreInfo(w, r, strings.Join(form["medID"], ""))
	} else {

		pageData := Medicines{}

		var medicines Medicines
		medicinesData, _ := os.Open("downloads/medicines.xml")
		defer medicinesData.Close()
		byteValue, _ := ioutil.ReadAll(medicinesData)
		xml.Unmarshal(byteValue, &medicines)

		pageData.LastUpdate = medicines.LastUpdate

		switch v := strings.Join(form["ean_code"], ""); v != "" {
		case true:
			for i := 0; i < len(medicines.Medicines); i++ {
				for j := 0; j < len(medicines.Medicines[i].Packages); j++ {
					switch strings.Contains(medicines.Medicines[i].Packages[j].EanCode, v) {
					case true:
						pageData.Medicines = append(pageData.Medicines, medicines.Medicines[i])
						break
					}
				}

				if len(pageData.Medicines) > 25 {
					break
				}
			}
			break
		}

		switch v := strings.Join(form["name"], ""); v != "" {
		case true:
			for i := 0; i < len(medicines.Medicines); i++ {
				switch strings.Contains(medicines.Medicines[i].ShortName, v) {
				case true:
					pageData.Medicines = append(pageData.Medicines, medicines.Medicines[i])
					break
				}

				switch strings.Contains(medicines.Medicines[i].FullName, v) {
				case true:
					pageData.Medicines = append(pageData.Medicines, medicines.Medicines[i])
					break
				}

				if len(pageData.Medicines) > 25 {
					break
				}
			}
			break
		}

		t, _ := template.ParseFiles("templates/medicines/results.html")
		t.Execute(w, pageData)
	}
}

func Page(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		Results(w, r)
	} else {
		t, _ := template.ParseFiles("templates/medicines/start.html")
		t.Execute(w, nil)
	}
}
