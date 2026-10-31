// SPDX-License-Identifier: MIT

package iec

import "testing"

// Whatever ParseRanges accepts is a set of ordered ranges within max, each of
// whose bounds the set contains.
func FuzzParseRanges(f *testing.F) {
	for _, s := range []string{"", "100,200-203", " 5, 0x10-18 ,,7 ", "9-5", "0xffffffff", "4294967296", "-", "1-", "a", "0b11", "0o7", "1_0"} {
		f.Add(s, uint32(MaxIOA))
	}
	f.Add("1-5", uint32(3))
	f.Add("0", uint32(0))
	f.Fuzz(func(t *testing.T, s string, max uint32) {
		r, err := ParseRanges(s, "number", max)
		if err != nil {
			return
		}
		for _, x := range r {
			if x.Lo > x.Hi || x.Hi > max {
				t.Fatalf("ParseRanges(%q, %d): range %+v", s, max, x)
			}
			if !r.Contains(x.Lo) || !r.Contains(x.Hi) {
				t.Fatalf("ParseRanges(%q, %d): %+v does not contain its bounds", s, max, x)
			}
		}
	})
}

// ParseIOAs bounds the list: a command cannot be made to expand a huge range.
func FuzzParseIOAs(f *testing.F) {
	for _, s := range []string{"", "1-200000", "1,2,3", "0-16777215", "16777216"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		r, err := ParseIOAs(s)
		if err != nil {
			return
		}
		if n := r.Count(); n > maxListed {
			t.Fatalf("ParseIOAs(%q): %d addresses, the maximum is %d", s, n, maxListed)
		}
	})
}

// A type is a number 1-255 or a known name, never type 0.
func FuzzParseTypes(f *testing.F) {
	for _, s := range []string{"", "M_SP_NA_1", "m_sp_na_1,C_SC_NA_1", "1", "0", "255", "256", "nope", ",,", " 1 , 2 "} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		types, err := ParseTypes(s)
		if err != nil {
			return
		}
		for id := range types {
			if id == 0 {
				t.Fatalf("ParseTypes(%q): type 0", s)
			}
		}
	})
}
