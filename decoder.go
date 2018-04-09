package csvparse

import (
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/blendlabs/go-name-parser"
)

// CSVDecoder holds the header field map and io.reader interface
type CSVDecoder struct {
	header map[string]int
	file   io.Reader
}

// NewDecoder allocates a new instance of CSVDecoder
func NewDecoder(input io.Reader) *CSVDecoder {
	return &CSVDecoder{
		header: map[string]int{},
		file:   input,
	}
}

// DecodeCSV unmarshalls CSV file to a specified struct type
func (d *CSVDecoder) DecodeCSV(v interface{}) error {
	// Optional timer function for determining function duration
	defer timeTrack(time.Now(), "DecodeCSVtoStruct")

	// getCSVRows grabs [][]strings from spcified input
	csvRows, err := GetCSVRows(d.file)
	if err != nil {
		return err
	}

	// returns an error if an empty file was provided
	if len(csvRows) == 0 {
		return errors.New("empty csv file")
	}
	// grab Header row
	headerRow := csvRows[0]
	// grab get remaining rows as body
	body := csvRows[1:]

	// check header row for duplicate fields
	// if duplicate fields found, return error
	if err := CheckForDoubleHeaderNames(headerRow); err != nil {
		return err
	}

	// check interface type (v)
	// if type is not a slice, return error
	slice := CheckInterfaceValue(v)
	if slice.Kind() != reflect.Slice {
		return errors.New("Only slice types are permited")
	}

	// check inner slice type
	// if type is not a struct, return error
	innerType := GetInnerSliceType(v)
	if innerType.Kind() != reflect.Struct {
		return errors.New("Only slice of stucts type permited")
	}

	// allocate new reflect stuct instance from inner type
	innerValueHdr := reflect.New(innerType)

	// range over header row
	for i, csvColumnHdr := range headerRow {
		// range over struct fields
		for j := 0; j < innerValueHdr.Elem().NumField(); j++ {
			var regex string
			// grab regex string from struct tag
			// if not struct tags provided, struct name will be used
			if rgx, ok := reflect.Indirect(innerValueHdr).Type().Field(j).Tag.Lookup("csv"); ok {
				switch rgx {
				case "-":
					regex = reflect.Indirect(innerValueHdr).Type().Field(j).Name
				default:
					regex = rgx
				}
			} else {
				continue
			}
			if regexp.MustCompile(regex).MatchString(csvColumnHdr) {
				d.header[reflect.Indirect(innerValueHdr).Type().Field(j).Name] = i
			}
		}
	}

	for _, csvRow := range body {
		innerValueRow := reflect.New(innerType)
		for j := 0; j < innerValueRow.Elem().NumField(); j++ {
			sFName := reflect.Indirect(innerValueRow).Type().Field(j).Name
			switch innerValueRow.Elem().Type().Field(j).Type {
			case reflect.TypeOf(time.Now()):
				if _, ok := d.header[sFName]; ok {
					val := ParseDate(csvRow[d.header[sFName]])
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			default:
				if _, ok := d.header[sFName]; ok {
					if format, ok := reflect.Indirect(innerValueRow).Type().Field(j).Tag.Lookup("fmt"); ok {
						if format != "-" {
							val := csvRow[d.header[sFName]]
							fmtvalue, err := FormatStringVals(format, val)
							if err != nil {
								log.Fatalln(err)
							}
							innerValueRow.Elem().FieldByName(sFName).SetString(fmtvalue)
						} else {
							val := csvRow[d.header[sFName]]
							innerValueRow.Elem().FieldByName(sFName).SetString(val)
						}
					}
				}
			}
			if j == innerValueRow.Elem().NumField()-1 {
				FullN := innerValueRow.Elem().FieldByName("Fullname")
				FN := innerValueRow.Elem().FieldByName("Firstname")
				LN := innerValueRow.Elem().FieldByName("Lastname")
				if FullN.String() != "" && (FN.String() == "" || LN.String() == "") {
					name := names.Parse(FullN.String())
					innerValueRow.Elem().FieldByName("Firstname").SetString(name.FirstName)
					innerValueRow.Elem().FieldByName("MI").SetString(name.MiddleName)
					innerValueRow.Elem().FieldByName("Lastname").SetString(name.LastName)
				}
				if fmt.Sprint(innerValueRow.Elem().FieldByName("Zip")) != "" {
					zip, zip4 := ParseZip(fmt.Sprint(innerValueRow.Elem().FieldByName("Zip")))
					innerValueRow.Elem().FieldByName("Zip").SetString(zip)
					if innerValueRow.Elem().FieldByName("Zp4").String() != "" {
						innerValueRow.Elem().FieldByName("Zip4").SetString(zip4)
					}
				}
			}
		}
		slice.Set(reflect.Append(slice, innerValueRow))
	}
	return nil
}
