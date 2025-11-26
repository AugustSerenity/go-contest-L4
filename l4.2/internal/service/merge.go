package service

import "sort"

func MergeAll(lines []string) []string {
	sort.Strings(lines)
	return lines
}
