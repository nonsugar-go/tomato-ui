package checkpoint

import (
	"fmt"
	"strings"
)

func buildKV(key, value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf(" %s \"%s\"", key, value)
}

func buildIndexedKV(key string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, v := range values {
		sb.WriteString(buildKV(fmt.Sprintf("%s.%d", key, i+1), v))
	}
	return sb.String()
}

func buildComment(desc string) string {
	return buildKV("comments", desc)
}

func buildTags(tags []string) string {
	return buildIndexedKV("tags", tags)
}

func mapStrings(in []string, m map[string]string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if mv, ok := m[v]; ok {
			v = mv
		}
		out = append(out, v)
	}
	return out
}
