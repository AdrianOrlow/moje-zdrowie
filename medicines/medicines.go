package medicines

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Medicines struct {
	XMLName    xml.Name `xml:"produktyLecznicze"`
	LastUpdate string   `xml:"stanNaDzien,attr"`
	Min        int
	Max        int
	PageNumber int
	Medicines  []Medicine `xml:"produktLeczniczy"`
}

type Medicine struct {
	XMLName          xml.Name  `xml:"produktLeczniczy"`
	ProductName      string    `xml:"nazwaProduktu,attr"`
	CommonName       string    `xml:"nazwaPowszechnieStosowana,attr"`
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

var medicines Medicines

func GetName(ean string) (name string, ID string) {
	medicinesData, _ := os.Open("downloads/medicines.xml")
	defer medicinesData.Close()
	byteValue, _ := ioutil.ReadAll(medicinesData)
	xml.Unmarshal(byteValue, &medicines)

	for i := 0; i < len(medicines.Medicines); i++ {
		for j := 0; j < len(medicines.Medicines[i].Packages); j++ {
			switch medicines.Medicines[i].Packages[j].EanCode == ean {
			case true:
				name = medicines.Medicines[i].ProductName
				ID = medicines.Medicines[i].ID
				break
			}
		}
	}
	return
}

func MoreInfo(w http.ResponseWriter, r *http.Request, medID string) {
	infoData := Medicine{}

	for i := 0; i < len(medicines.Medicines); i++ {
		for j := 0; j < len(medicines.Medicines[i].Packages); j++ {
			switch medicines.Medicines[i].ID == medID {
			case true:
				infoData = medicines.Medicines[i]
				break
			}
		}
	}

	t, _ := template.ParseFiles("templates/medicines/info.html")
	t.Execute(w, infoData)
}

func CheckEanCode(w http.ResponseWriter, r *http.Request, ean string) {
	for i := 0; i < len(medicines.Medicines); i++ {
		for j := 0; j < len(medicines.Medicines[i].Packages); j++ {
			switch medicines.Medicines[i].Packages[j].EanCode == ean {
			case true:
				MoreInfo(w, r, medicines.Medicines[i].ID)
			}
		}
	}
}

func CheckName(min, max int, name string) Medicines {
	data := Medicines{}
	for i := min; i < len(medicines.Medicines); i++ {
		switch strings.Contains(strings.ToLower(medicines.Medicines[i].CommonName), strings.ToLower(name)) {
		case true:
			data.Medicines = append(data.Medicines, medicines.Medicines[i])
			break
		}

		switch strings.Contains(strings.ToLower(medicines.Medicines[i].ProductName), strings.ToLower(name)) {
		case true:
			data.Medicines = append(data.Medicines, medicines.Medicines[i])
			break
		}

		if len(data.Medicines) > max {
			data.Min = i
			data.Max = max
			data.PageNumber = int(float64(max) / 25.0)
			break
		}
	}
	data.LastUpdate = medicines.LastUpdate
	return data
}

func Results(w http.ResponseWriter, r *http.Request, min, max int) {
	r.ParseForm()
	form := r.Form

	if strings.Join(form["medID"], "") != "" {
		MoreInfo(w, r, strings.Join(form["medID"], ""))
	} else {
		switch ean := strings.Join(form["ean_code"], ""); ean != "" {
		case true:
			CheckEanCode(w, r, ean)
			return
		}

		var pageData Medicines

		switch name := strings.Join(form["name"], ""); name != "" {
		case true:
			pageData = CheckName(min, max, name)
			break
		}

		t, _ := template.ParseFiles("templates/medicines/results.html")
		t.Execute(w, pageData)
	}
}

func Page(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		medicinesData, _ := os.Open("downloads/medicines.xml")
		defer medicinesData.Close()
		byteValue, _ := ioutil.ReadAll(medicinesData)
		xml.Unmarshal(byteValue, &medicines)

		var min, max int = 0, 25
		if len(r.Form["min"]) > 0 {
			max, _ = strconv.Atoi(r.FormValue("max")[0:])
			min, _ = strconv.Atoi(r.FormValue("min")[0:])
		}
		Results(w, r, min, max)
	} else {
		t, _ := template.ParseFiles("templates/medicines/start.html")
		t.Execute(w, nil)
	}
}
