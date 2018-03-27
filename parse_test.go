package csvparse_test

import (
	"testing"
	"time"

	cp "github.com/rssenar/csvparse"
)

func Test_tCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{" apPlE ", "Apple"},
		{" sUn ", "Sun"},
		{" nOaH  ", "Noah"},
	}
	for _, c := range cases {
		out := cp.TCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_UCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{" aPPle  ", "APPLE"},
		{"  sUN ", "SUN"},
		{" noAH ", "NOAH"},
	}
	for _, c := range cases {
		out := cp.UCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_LCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"   APPLE  ", "apple"},
		{" SUN ", "sun"},
		{"  NOAH   ", "noah"},
	}
	for _, c := range cases {
		out := cp.LCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_parseZip(t *testing.T) {
	cases := []struct {
		input, zip, zip4 string
	}{
		{"92882-1234", "92882", "1234"},
		{"928821234", "92882", "1234"},
		{"928821234", "92882", "1234"},
		{"9288212", "9288212", ""},
		{"92882123456", "92882123456", ""},
	}
	for _, c := range cases {
		zip, zip4 := cp.ParseZip(c.input)
		if zip != c.zip {
			t.Errorf("Zip should be %v but got %v", c.zip, zip)
		}
		if zip4 != c.zip4 {
			t.Errorf("Zip4 should be %v but got %v", c.zip4, zip4)
		}
	}
}

func Test_FormatPhone(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"9493237895", "(949) 323-7895"},
		{"3237895", "323-7895"},
		{"94932", ""},
		{"94932456748912", ""},
	}
	for _, c := range cases {
		out := cp.FormatPhone(c.input)
		if out != c.expected {
			t.Errorf("Phone should be %v but got %v", c.expected, out)
		}
	}
}

func Test_StripSep(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"#$*string&()&", "string"},
		{"#   $*   string   &()&   ", "string"},
	}
	for _, c := range cases {
		out := cp.StripSep(c.input)
		if out != c.expected {
			t.Errorf("Phone should be %v but got %v", c.expected, out)
		}
	}
}

func Test_TrimZeros(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"00000123", "123"},
		{"000000000000000000123", "123"},
		{"0123", "123"},
	}
	for _, c := range cases {
		out := cp.TrimZeros(c.input)
		if out != c.expected {
			t.Errorf("number should be %v but got %v", c.expected, out)
		}
	}
}

func Test_ParseDate(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"12/31/2003", extDate("12/31/2003")},
		{"12-31-2003", extDate("12-31-2003")},
		{"1-3-03", extDate("1-3-03")},
		{"", extDate("")},
	}
	for _, c := range cases {
		out := cp.ParseDate(c.input)
		if out != c.expected {
			t.Errorf("Date should be %v but got %v", c.expected, out)
		}
	}
}

func extDate(date string) time.Time {
	formats := []string{"1/2/2006", "1-2-2006", "1/2/06", "1-2-06", "2006/1/2", "2006-1-2", time.RFC3339}
	for _, f := range formats {
		if date, err := time.Parse(f, date); err == nil {
			return date
		}
	}
	return time.Time{}
}
