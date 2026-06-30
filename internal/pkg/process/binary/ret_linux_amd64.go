// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build 386 || amd64

package binary

import (
	"golang.org/x/arch/x86/x86asm"
)

func findRetInstructions(data []byte) ([]uint64, error) {
	n := uint64(len(data))

	var returnOffsets []uint64
	var index uint64
	for index < n {
		instruction, err := x86asm.Decode(data[index:], 64)
		if err != nil {
			// Keep scanning when x86asm does not recognize a newer instruction.
			index++
			continue
		}

		if instruction.Op == x86asm.RET {
			returnOffsets = append(returnOffsets, index)
		}

		index += uint64(max(1, instruction.Len)) //nolint:gosec  // Underflow handled.
	}

	return returnOffsets, nil
}
