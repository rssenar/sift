package csvparse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// GetCSVRows get [][]strins from io.reader
func GetCSVRows(r io.Reader) ([][]string, error) {
	rdr := csv.NewReader(r)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%v : unable to read file", err)
	}
	return rows, nil
}

// GetInnerSliceType gets inner slice tyoe with reflect and return non-pointer value
func GetInnerSliceType(v interface{}) reflect.Type {
	outType := reflect.TypeOf(v)
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}
	inType := outType.Elem()
	if inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	return inType
}

// CheckInterfaceValue get interface tyoe with reflect and return non-pointer value
func CheckInterfaceValue(v interface{}) reflect.Value {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

// CheckForDoubleHeaderNames checks for duplicate header fields
// returns error if dupes found
func CheckForDoubleHeaderNames(hdrs []string) error {
	headerMap := make(map[string]bool, len(hdrs))
	for _, v := range hdrs {
		if _, ok := headerMap[v]; ok {
			return fmt.Errorf("Repeated header name: %v", v)
		}
		headerMap[v] = true
	}
	return nil
}

// FormatStringVals applies formating to string
// based on "fmt" struct tag
func FormatStringVals(format, val string) (string, error) {
	switch format {
	case "tc":
		return TCase(val), nil
	case "uc":
		return UCase(val), nil
	case "lc":
		return LCase(val), nil
	case "fp":
		return FormatPhone(val), nil
	case "ss":
		return StripSep(val), nil
	default:
		return "", errors.New("Invalid string format: use [tc, uc, lc, fp, ss]")
	}
}

// TCase transforms string to title case and trims leading & trailing white space
func TCase(f string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(f)))
}

// UCase transforms string to upper case and trims leading & trailing white space
func UCase(f string) string {
	return strings.TrimSpace(strings.ToUpper(f))
}

// LCase transforms string to lower case and trims leading & trailing white space
func LCase(f string) string {
	return strings.TrimSpace(strings.ToLower(f))
}

// ParseZip perses ZIP code to Zip & Zip4
func ParseZip(zip string) (string, string) {
	if zip == "" {
		return "", ""
	}
	switch {
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]$`).MatchString(zip):
		return TrimZeros(zip[:5]), TrimZeros(zip[5:])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, "-")
		return TrimZeros(zsplit[0]), TrimZeros(zsplit[1])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9] [0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, " ")
		return TrimZeros(zsplit[0]), TrimZeros(zsplit[1])
	default:
		return zip, ""
	}
}

// TrimZeros removed leading Zeros
func TrimZeros(s string) string {
	l := len(s)
	for i := 1; i <= l; i++ {
		s = strings.TrimPrefix(s, "0")
	}
	return s
}

// FormatPhone re-formats phone field
func FormatPhone(p string) string {
	p = StripSep(p)
	switch len(p) {
	case 10:
		return fmt.Sprintf("(%v) %v-%v", p[0:3], p[3:6], p[6:10])
	case 7:
		return fmt.Sprintf("%v-%v", p[0:3], p[3:7])
	default:
		return ""
	}
}

// StripSep removes irrelevant characters
func StripSep(p string) string {
	sep := []string{"'", "#", "%", "$", "-", "+", ".", "*", "(", ")", ":", ";", "{", "}", "|", "&", " "}
	for _, v := range sep {
		p = strings.Replace(p, v, "", -1)
	}
	return p
}

// ParseDate converts date string to time.Time
func ParseDate(d string) time.Time {
	if d != "" {
		formats := []string{"1/2/2006", "1-2-2006", "1/2/06", "1-2-06", "2006/1/2", "2006-1-2", time.RFC3339}
		for _, f := range formats {
			if date, err := time.Parse(f, d); err == nil {
				return date
			}
		}
	}
	return time.Time{}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}