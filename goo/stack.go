package goo

import (
	"bytes"
	"fmt"
	"strings"
)

func locatePanic(stack []byte) []string {
	lines := bytes.Split(stack, []byte("\n"))
	if len(lines) < 3 {
		return []string{string(stack)}
	}

	lines = lines[1:]
	var stacks []string
	for i := 0; i < len(lines); i += 2 {
		if i < len(lines)-1 {
			stacks = append(stacks, fmt.Sprintf("(%s::%s)", lines[i], bytes.TrimSpace(lines[i+1])))
		}
	}

	for i, stack := range stacks {
		if strings.Contains(stack, "src/runtime/panic.go:") && i < len(stacks)-1 {
			return stacks[i+1:]
		}
	}

	return stacks
}
