// SPDX-License-Identifier: MIT

package station

import "testing"

// FuzzParse feeds model files: the operator's, or one handed to them. A model
// that parses and validates must survive its own encoding.
func FuzzParse(f *testing.F) {
	data, err := Demo().Encode(false)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(data)
	if data, err = Demo().Encode(true); err != nil {
		f.Fatal(err)
	}
	f.Add(data)
	for _, s := range []string{"", "{}", "null", "- a\n- b\n", "station: {common_address: 1}\npoints:\n  - {ioa: 1, type: M_SP_NA_1, value: true}\n", "station: {common_address: 1, colour: red}", "points: [{}]\n", "\xff\xfe", "a: &a [*a]\n"} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		m, err := Parse(data)
		if err != nil {
			return
		}
		// Validate reads the files a model names: no input gets to pick them here.
		m.Files = nil
		if m.Validate() != nil {
			return
		}
		out, err := m.Encode(false)
		if err != nil {
			t.Fatalf("a valid model does not encode: %v\n%s", err, data)
		}
		back, err := Parse(out)
		if err != nil {
			t.Fatalf("a valid model does not parse back: %v\n%s", err, out)
		}
		if err := back.Validate(); err != nil {
			t.Fatalf("a valid model does not validate back: %v\n%s", err, out)
		}
	})
}
