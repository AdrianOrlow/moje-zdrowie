package pharmacies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jasonwinn/geocoder"
)

type Pharmacies struct {
	XMLName  xml.Name   `xml:"Apteki"`
	Pharmacy []Pharmacy `xml:"Apteka"`
}

type Pharmacy struct {
	XMLName           xml.Name     `xml:"Apteka"`
	ID                string       `xml:"id,attr"`
	Name              string       `xml:"nazwa,attr"`
	Type              string       `xml:"rodzaj,attr"`
	Phone             string       `xml:"numerTelefonu,attr"`
	Status            string       `xml:"status,attr"`
	TemporarilyClosed string       `xml:"czasowoNieczynna,attr"`
	OpeningDate       string       `xml:"dataUruchomieniaApteki,attr"`
	Website           string       `xml:"adresWWWApteki,attr"`
	Shop              string       `xml:"adresWWWSprzedazyWysylkowej,attr"`
	RangeOfActivites  string       `xml:"zakresDzialalnosci,attr"`
	Lat               float64      `xml:"szerokoscGeograficzna,attr"`
	Lng               float64      `xml:"dlugoscGeograficzna,attr"`
	Owners            []Owner      `xml:"Wlasciciele>Wlasciciel"`
	Address           Address      `xml:"Adres"`
	Manager           Manager      `xml:"KierownikApteki"`
	OpeningDays       []OpeningDay `xml:"DniPracyApteki>DniPracy>DzienPracy"`
}

type Owner struct {
	LegalForm string `xml:"formaPrawna,attr"`
	FirstName string `xml:"imie,attr"`
	KRS       string `xml:"krs,attr"`
	Name      string `xml:"nazwa,attr"`
	Surname   string `xml:"nazwisko,attr"`
	NIP       string `xml:"nip,attr"`
	REGON     string `xml:"regon,attr"`
}

type Address struct {
	ZIPCode         string `xml:"kodPocztowy,attr"`
	City            string `xml:"miejscowosc,attr"`
	StreetType      string `xml:"typUlicy,attr"`
	StreetName      string `xml:"nazwaUlicy,attr"`
	HouseNumber     string `xml:"numerDomu,attr"`
	ApartmentNumber string `xml:"numerLokalu,attr"`
	Voivodeship     string `xml:"wojewodztwo,attr"`
}

type Manager struct {
	FirstName string `xml:"imie,attr"`
	Surname   string `xml:"nazwisko,attr"`
}

type OpeningDay struct {
	DayName   string `xml:"dzienTygodnia,attr"`
	OpenFrom  string `xml:"otwartaOd,attr"`
	OpenUntil string `xml:"otwartaDo,attr"`
	AllDay    string `xml:"calodobowa,attr"`
}

/* Markers */

type PharmacyMarkers struct {
	PharmacyMarkers []PharmacyMarker
}

type PharmacyMarker struct {
	ID    string
	Image string
	Lat   string
	Lng   string
}

func ChooseImage(Type string) (link string) {
	link = "static/img/pharmacies/"
	switch Type {
	case "apteka ogólnodostępna":
		link += ""
		break
	case "punkt apteczny":
		link += ""
		break
	case "dział farmacji szpitalnej":
		link += ""
		break
	default:
		link += ""
		break
	}

	return
}

func GetCoord(a string) (string, string) {
	lat, lng, _ := geocoder.Geocode(a)
	return fmt.Sprint(lat), fmt.Sprint(lng)
}

func (a Address) addr() string {
	return a.StreetType + " " + a.StreetName + " " + a.HouseNumber + ", " + a.City + " " + a.ZIPCode + ", " + a.Voivodeship + ", Poland"
}

func ConvertPharmaciesToMarkers() {
	var pharmacyMarkers PharmacyMarkers
	var pharmacies Pharmacies

	pharmaciesData, _ := os.Open("downloads/pharmacies.xml")
	defer pharmaciesData.Close()
	byteValue, _ := ioutil.ReadAll(pharmaciesData)
	xml.Unmarshal(byteValue, &pharmacies)

	geocoder.SetAPIKey("j4wlV4OuOLNcT3ZkAGo7h5qxjtkrsAL7")
	for j, i := range pharmacies.Pharmacy {
		pM := PharmacyMarker{}

		addr := i.Address.addr()
		pM.Lat, pM.Lng = GetCoord(addr)

		fmt.Println(j, "/", len(pharmacies.Pharmacy))

		pM.ID = i.ID
		pM.Image = ChooseImage(i.Type)

		pharmacyMarkers.PharmacyMarkers = append(pharmacyMarkers.PharmacyMarkers, pM)
	}

	markers, _ := json.Marshal(pharmacyMarkers)
	ioutil.WriteFile("static/json/pharmacies-markers.json", markers, 0644)
}
