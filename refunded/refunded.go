package refunded

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"moje-zdrowie/medicines"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Refunded struct {
	XMLName           xml.Name           `xml:"lekiRefundowane"`
	RefundedMedicines []RefundedMedicine `xml:"lekRefundowany"`
}

type RefundedMedicine struct {
	XMLName         xml.Name `xml:"lekRefundowany"`
	EanCode         string   `xml:"kodEan,attr"`
	MedicineAmount  float64  `xml:"iloscWydanegoLeku,attr"`
	RefundAmount    float64  `xml:"kwotaRefundacji,attr"`
	RefundByPatient string
	ProductName     string
	ID              string
}

var refunded Refunded

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func Results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	refundedMed := RefundedMedicine{}

	switch ean := strings.Join(form["ean_code"], ""); ean != "" {
	case true:
		for i := 0; i < len(refunded.RefundedMedicines); i++ {
			switch refunded.RefundedMedicines[i].EanCode == ean {
			case true:
				refundedMed = refunded.RefundedMedicines[i]
				break
			}
		}
		refundedMed.RefundByPatient = FloatToString(refundedMed.RefundAmount / refundedMed.MedicineAmount)
		refundedMed.ProductName, refundedMed.ID = medicines.GetName(ean)
		break
	}

	t, _ := template.ParseFiles("templates/refunded/results.html")
	t.Execute(w, refundedMed)
}

func Page(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		refundedData, _ := os.Open("downloads/refunded.xml")
		defer refundedData.Close()
		byteValue, _ := ioutil.ReadAll(refundedData)
		xml.Unmarshal(byteValue, &refunded)

		Results(w, r)
	} else {
		t, _ := template.ParseFiles("templates/refunded/start.html")
		t.Execute(w, nil)
	}
}
