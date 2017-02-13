package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// box is generated from an XSD element
type box struct {
	ID        string    `xml:"id"`
	Type      string    `xml:"type"`
	Subtype   int       `xml:"subtype"`
	Name      name      `xml:"name"`
	Ico       string    `xml:"ico"`
	Address   address   `xml:"address"`
	Pdz       bool      `xml:"pdz"`
	Ovm       bool      `xml:"ovm"`
	Hierarchy hierarchy `xml:"hierarchy"`
}

// name is generated from an XSD element
type name struct {
	Person    person `xml:"person"`
	TradeName string `xml:"tradeName"`
}

// person is generated from an XSD element
type person struct {
	FirstName       string `xml:"firstName"`
	LastName        string `xml:"lastName"`
	MiddleName      string `xml:"middleName"`
	LastNameAtBirth string `xml:"lastNameAtBirth"`
}

// address is generated from an XSD element
type address struct {
	Code         int    `xml:"code"`
	City         string `xml:"city"`
	District     string `xml:"district"`
	Street       string `xml:"street"`
	Cp           string `xml:"cp"`
	Co           string `xml:"co"`
	Ce           string `xml:"ce"`
	Zip          string `xml:"zip"`
	Region       string `xml:"region"`
	AddressPoint int    `xml:"addressPoint"`
	State        string `xml:"state"`
	FullAddress  string `xml:"fullAddress"`
}

// hierarchy is generated from an XSD element
type hierarchy struct {
	IsMaster bool   `xml:"isMaster"`
	MasterID string `xml:"masterId"`
}

// TODO: generuj header (zase reflect)
// TODO: gzip pri psani
// TODO: panic -> error
func main() {
	dr := "data"
	tdr := "csv"
	fns, _ := filepath.Glob(filepath.Join(dr, "*.xml.gz"))

	for _, fn := range fns {
		fmt.Printf("Parsuju %s\n", fn)

		_, cfn := filepath.Split(fn)
		tfn := filepath.Join(tdr, fmt.Sprintf("%s.csv", cfn[:strings.Index(cfn, ".")]))

		konvertuj(fn, tfn)
	}
}

// rekurzivne konvertuj struct v []string
func exportuj(b interface{}) []string {
	dt := make([]string, 0)

	tp := reflect.TypeOf(b)
	if tp.Kind() != reflect.Struct {
		panic("expected a struct")
	}
	vl := reflect.ValueOf(b)
	for j := 0; j < tp.NumField(); j++ {
		switch tp.Field(j).Type.Kind() {
		case reflect.String:
			dt = append(dt, vl.Field(j).String())

		case reflect.Bool:
			dt = append(dt, strconv.FormatBool(vl.Field(j).Bool()))

		case reflect.Int:
			dt = append(dt, strconv.FormatInt(vl.Field(j).Int(), 10))
		default:
			// jestli jdem do nestovanyho pole, musi to byt struct
			if tp.Field(j).Type.Kind() != reflect.Struct {
				panic(tp.Field(j))
			}
			// nestovana data
			ndt := exportuj(vl.Field(j).Interface())
			for _, v := range ndt {
				dt = append(dt, v)
			}
		}
	}

	return dt
}

func konvertuj(fn string, tfn string) {
	fl, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer fl.Close()

	fg, err := gzip.NewReader(fl)
	if err != nil {
		panic(err)
	}
	xd := xml.NewDecoder(fg)

	tfl, err := os.Create(tfn)
	if err != nil {
		panic(err)
	}
	defer tfl.Close()
	tcw := csv.NewWriter(tfl)

	tot := 0
	t := time.Now()
	for {
		tk, err := xd.Token()
		if err != nil {
			break
		}
		if reflect.TypeOf(tk).String() != "xml.StartElement" {
			continue
		}
		el := tk.(xml.StartElement)
		if !(el.Name.Local == "box") {
			continue
		}
		var b box
		if err := xd.DecodeElement(&b, &el); err != nil {
			panic(err)
		}

		dt := exportuj(b)
		if err := tcw.Write(dt); err != nil {
			panic(err)
		}

		tot += 1
	}
	tcw.Flush()
	fmt.Println(tot, time.Since(t))
}
