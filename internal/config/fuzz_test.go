// SPDX-License-Identifier: MIT

package config

import "testing"

// A target comes from --url, a flag or a script. What ParseURL accepts must be
// usable: a host and a port.
func FuzzParseURL(f *testing.F) {
	for _, s := range []string{
		"tcp://10.0.0.1:2404", "tls://host", "tcp://[::1]:2404", "TCP://host:1", "tcp://", "tcp://host:0",
		"tcp://host:65536", "tcp://u@host", "tcp://host/path", "ftp://host", "", " tcp://host ", "%zz",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		tgt, err := ParseURL(raw)
		if err != nil {
			if tgt != (Target{}) {
				t.Fatalf("ParseURL(%q): %+v with error %v", raw, tgt, err)
			}
			return
		}
		if tgt.Host == "" || tgt.Port == 0 {
			t.Fatalf("ParseURL(%q): %+v", raw, tgt)
		}
	})
}
