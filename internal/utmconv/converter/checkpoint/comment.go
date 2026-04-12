package checkpoint

import "fmt"

func buildComment(desc string) string {
	if desc == "" {
		return ""
	}
	return fmt.Sprintf(" comments \"%s\"", desc)
}

func buildTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	tagsStr := " "
	for i, tag := range tags {
		tagsStr += fmt.Sprintf("tags.%d \"%s\"", i, tag)
	}
	return tagsStr
}
