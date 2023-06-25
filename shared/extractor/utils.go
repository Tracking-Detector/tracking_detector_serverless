package extractor

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
