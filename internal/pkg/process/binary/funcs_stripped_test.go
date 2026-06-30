// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package binary

import (
	"debug/elf"
	"encoding/binary"
	"testing"
)

func TestRuntimeTextFromPclntabFallback(t *testing.T) {
	text := &elf.Section{
		SectionHeader: elf.SectionHeader{Addr: 0x4d50200},
	}

	tests := []struct {
		name string
		pcln []byte
		want uint64
	}{
		{
			name: "zero 32-bit text start falls back to text section",
			pcln: newPclntab(4, 0),
			want: text.Addr,
		},
		{
			name: "zero 64-bit text start falls back to text section",
			pcln: newPclntab(8, 0),
			want: text.Addr,
		},
		{
			name: "non-zero 32-bit text start is preserved",
			pcln: newPclntab(4, 0x401000),
			want: 0x401000,
		},
		{
			name: "non-zero 64-bit text start is preserved",
			pcln: newPclntab(8, 0x401000),
			want: 0x401000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtimeTextFromPclntab(tt.pcln, text)
			if err != nil {
				t.Fatalf("runtimeTextFromPclntab returned error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("runtimeTextFromPclntab() = %#x, want %#x", got, tt.want)
			}
		})
	}
}

func TestRuntimeTextFromPclntabZeroWithoutTextSection(t *testing.T) {
	got, err := runtimeTextFromPclntab(newPclntab(8, 0), nil)
	if err != nil {
		t.Fatalf("runtimeTextFromPclntab returned error: %v", err)
	}

	if got != 0 {
		t.Fatalf("runtimeTextFromPclntab() = %#x, want 0", got)
	}
}

func TestRuntimeTextFromPclntabInvalidPointerSize(t *testing.T) {
	_, err := runtimeTextFromPclntab(newPclntab(16, 0), nil)
	if err == nil {
		t.Fatal("runtimeTextFromPclntab returned nil error")
	}
}

func newPclntab(ptrSize byte, runtimeText uint64) []byte {
	pcln := make([]byte, 128)
	pcln[7] = ptrSize

	off := 8 + 2*int(ptrSize)
	switch ptrSize {
	case 4:
		binary.LittleEndian.PutUint32(pcln[off:], uint32(runtimeText))
	case 8:
		binary.LittleEndian.PutUint64(pcln[off:], runtimeText)
	}

	return pcln
}
