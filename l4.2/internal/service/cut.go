package service

import (
	"strings"
)

type Options struct {
	Fields    string
	Delimiter string
	Separated bool

	FieldsSet map[int]bool
}

func (s *Service) CutLines(lines []string) []string {
	if s.Opts.FieldsSet == nil {
		s.Opts.FieldsSet = parseFields(s.Opts.Fields)
	}

	var result []string
	for _, line := range lines {
		out := processLine(line, s.Opts)
		if out != "" {
			result = append(result, out)
		}
	}
	return result
}

func processLine(line string, opts Options) string {
	if opts.Separated && !strings.Contains(line, opts.Delimiter) {
		return ""
	}

	parts := strings.Split(line, opts.Delimiter)

	if len(opts.FieldsSet) == 0 {
		return line
	}

	var selected []string
	for i := range parts {
		if opts.FieldsSet[i+1] {
			selected = append(selected, parts[i])
		}
	}

	return strings.Join(selected, opts.Delimiter)
}

func parseFields(spec string) map[int]bool {
	fields := make(map[int]bool)
	if spec == "" {
		return fields
	}
	for _, part := range strings.Split(spec, ",") {
		if strings.Contains(part, "-") {
			b := strings.Split(part, "-")
			start := atoi(b[0])
			end := atoi(b[1])
			for i := start; i <= end; i++ {
				fields[i] = true
			}
		} else {
			fields[atoi(part)] = true
		}
	}
	return fields
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			n = n*10 + int(r-'0')
		}
	}
	return n
}
