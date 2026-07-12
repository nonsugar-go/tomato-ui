package fortigate

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

type replaceRule struct {
	match   *regexp.Regexp
	replace *regexp.Regexp
	repl    []byte
	global  bool
}

func (r replaceRule) Apply(line []byte) []byte {
	if !r.match.Match(line) {
		return line
	}

	if r.global {
		return r.replace.ReplaceAll(line, r.repl)
	}

	return replaceFirst(line, r.replace, r.repl)
}

const fgSeqKeysPattern = `(member|gui-vdom-menu-favorites|vci-string|internet-service-name|service|srcaddr|dstaddr)`

var fgRules = []replaceRule{
	{
		regexp.MustCompile(`^\s*` + fgSeqKeysPattern + `: "[^"]+" "`),
		regexp.MustCompile(`" "`),
		[]byte(`", "`), true,
	},
	{
		regexp.MustCompile(`^\s*` + fgSeqKeysPattern + `: "[^"]+"`),
		regexp.MustCompile(`$`),
		[]byte(`]`), false,
	},
	{
		regexp.MustCompile(`^\s*` + fgSeqKeysPattern + `: "[^"]+"`),
		regexp.MustCompile(`: `),
		[]byte(`: [`), false,
	},

	// KEY: "a" "b"
	// {
	// 	regexp.MustCompile(`^\s*\S+: "[^"]+" "`),
	// 	regexp.MustCompile(`" "`),
	// 	[]byte(`", "`), true,
	// },
	// {
	// 	regexp.MustCompile(`^\s*\S+: "[^"]+", "`),
	// 	regexp.MustCompile(`$`),
	// 	[]byte(`]`), false,
	// },
	// {
	// 	regexp.MustCompile(`^\s*\S+: "[^"]+", "`),
	// 	regexp.MustCompile(`: `),
	// 	[]byte(`: [`), false,
	// },

	// - *.bat:
	{
		regexp.MustCompile(`^\s*- [*].+:$`),
		regexp.MustCompile(`:$`),
		[]byte(`':`), false,
	},
	{
		regexp.MustCompile(`^\s*- [*].+:$`),
		regexp.MustCompile(`- `),
		[]byte(`- '`), false,
	},

	// - 0087654:
	{
		regexp.MustCompile(`^\s*- 0.+:$`),
		regexp.MustCompile(`:$`),
		[]byte(`':`), false,
	},
	{
		regexp.MustCompile(`^\s*- 0.+:$`),
		regexp.MustCompile(`- `),
		[]byte(`- '`), false,
	},
}

func FgValidYaml(reader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 16<<20)
	for scanner.Scan() {
		line := scanner.Bytes()
		for _, r := range fgRules {
			line = r.Apply(line)
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), scanner.Err()
}

func replaceFirst(src []byte, re *regexp.Regexp, repl []byte) []byte {
	match := re.FindIndex(src)
	if match == nil {
		return src
	}
	dst := make([]byte, 0, len(src)-match[1]+match[0]+len(repl))
	dst = append(dst, src[:match[0]]...)
	dst = append(dst, repl...)
	dst = append(dst, src[match[1]:]...)
	return dst
}
