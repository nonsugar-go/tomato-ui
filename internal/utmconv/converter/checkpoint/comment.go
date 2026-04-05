package checkpoint

import "fmt"

func buildComment(desc string) string {
	if desc == "" {
		return ""
	}
	return fmt.Sprintf(" comments \"%s\"", desc)
}
