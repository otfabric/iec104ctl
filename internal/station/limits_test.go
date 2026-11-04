// SPDX-License-Identifier: MIT

package station

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otfabric/go-iec104/asdu"
	"github.com/otfabric/go-iec104/server"
)

// The length a directory reports fits the 24 bits of a length of file; what does not fit is
// reported as the largest length, never wrapped.
func TestFileLengthBounds(t *testing.T) {
	for n, want := range map[int]uint32{
		0:                      0,
		1:                      1,
		asdu.MaxFileLength - 1: asdu.MaxFileLength - 1,
		asdu.MaxFileLength:     asdu.MaxFileLength,
		asdu.MaxFileLength + 1: asdu.MaxFileLength,
		1 << 32:                asdu.MaxFileLength,
		1<<32 + 5:              asdu.MaxFileLength, // would wrap to 5 as a uint32
		-1:                     asdu.MaxFileLength,
	} {
		if got := fileLength(n); got != want {
			t.Errorf("fileLength(%d) = %d, want %d", n, got, want)
		}
	}
}

// The upload directory and what is delivered into it are the operator's: the directory is
// not open to others, a file not to group or others, whatever the umask.
func TestUploadsArePrivate(t *testing.T) {
	dir := t.TempDir()
	model := "station: {common_address: 3}\npoints: [{ioa: 1, type: M_SP_NA_1, value: true}]\n" +
		"uploads: {accept: true, dir: uploads/nested}\n"
	path := filepath.Join(dir, "station.yaml")
	if err := os.WriteFile(path, []byte(model), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	st, err := New(m, Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"uploads", "uploads/nested"} {
		info, err := os.Stat(filepath.Join(dir, d))
		if err != nil {
			t.Fatal(err)
		}
		if perm := info.Mode().Perm(); perm&0o007 != 0 {
			t.Errorf("%s has mode %v: others can access it", d, perm)
		}
	}
	if err := st.StoreFile(nil, 3, 77, 4, server.NewFile([]byte("delivered"), 4)); err != nil {
		t.Fatal(err)
	}
	stored := filepath.Join(dir, "uploads", "nested", "ioa-77-name-4.bin")
	info, err := os.Stat(stored)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm&0o077 != 0 {
		t.Errorf("a delivered file has mode %v: group or others can access it", perm)
	}
	if data, ok := st.Uploaded(77); !ok || string(data) != "delivered" {
		t.Errorf("Uploaded: %q, %v", data, ok)
	}
	// The directory listing reports the length of the delivered file.
	for _, e := range st.directory(nil, 3, 0) {
		if e.IOA == 77 && e.Length != 9 {
			t.Errorf("directory entry %+v", e)
		}
	}
}
