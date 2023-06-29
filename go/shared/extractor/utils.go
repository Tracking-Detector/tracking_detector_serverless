package extractor

import (
	"unicode/utf16"
	"unicode/utf8"
)

func GetTypes() []string {
	return []string{
		"xmlhttprequest",
		"image",
		"font",
		"script",
		"stylesheet",
		"ping",
		"sub_frame",
		"other",
		"main_frame",
		"csp_report",
		"object",
		"media"}
}

func GetFrameTypes() []string {
	return []string{"outermost_frame", "fenced_frame", "sub_frame"}
}

func GetMethods() []string {
	return []string{"GET",
		"POST",
		"OPTIONS",
		"HEAD",
		"PUT",
		"DELETE",
		"SEARCH",
		"PATCH"}
}

func URL_EXTRACTOR(s string) []int {
	encoded := make([]int, 200)
	count := 199
	for i := len(s) - 1; i >= 0; i-- {
		c, _ := utf8.DecodeRuneInString(s[i:i])
		utf16Enc := utf16.Encode([]rune{c})
		encoded[count] = (int(utf16Enc[0]) % 98) + 1
		if count == 0 {
			break
		}
		count--
	}
	return encoded
}

func FRAME_TYPE_EXTRACTOR(s string) []int {
	for i, val := range GetFrameTypes() {
		if val == s {
			return []int{i + 1}
		}
	}
	return []int{0}
}

func METHOD_EXTRACTOR(s string) []int {
	for i, val := range GetMethods() {
		if val == s {
			return []int{i + 1}
		}
	}
	return []int{0}
}

func TYPE_EXTRACTOR(s string) []int {
	for i, val := range GetTypes() {
		if val == s {
			return []int{i + 1}
		}
	}
	return []int{0}
}

func LABEL_EXTRACTOR(b bool) []int {
	if b {
		return []int{1}
	}
	return []int{0}
}

func REQUEST_HEADER_REFERER_EXTRACTOR(headers []map[string]string) []int {
	for _, header := range headers {
		if val, exists := header["name"]; exists && val == "Referer" {
			return []int{1}
		}
	}
	return []int{0}
}
