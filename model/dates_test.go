package model

import "testing"

func TestNormalizeBirthDate(t *testing.T) {
	cases := map[string]string{"41050": "21.05.2012", "41050.0": "21.05.2012", "01.02.2010": "01.02.2010", "": ""}
	for in, want := range cases {
		if got := NormalizeBirthDate(in); got != want {
			t.Errorf("%q: got %q want %q", in, got, want)
		}
	}
}
