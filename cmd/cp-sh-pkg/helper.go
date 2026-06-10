package main

import (
	"strconv"
	"strings"
)

func b2s(b bool) string {
	return strconv.FormatBool(b)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func join(ss []string) string {
	return strings.Join(ss, ";")
}

func joinNames(objs []CPObject) string {
	var sb strings.Builder
	for i, o := range objs {
		sb.WriteString(o.Name)
		if i != 0 {
			sb.WriteRune(';')
		}
	}
	return sb.String()
}
