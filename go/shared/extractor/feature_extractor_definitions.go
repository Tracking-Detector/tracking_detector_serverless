package extractor

var EXTRACTORS = InitExtractors()

func InitExtractors() []Extractor {
	// Create Extractor with dimensions [204,1]
	extractor204 := NewExtractor("Extractor204", "This Extractor extracts a feature vector of [204,1] it includes 200 chars of the URL, the request type, the request method, the request frame_type and the presence of the referrer header")
	extractor204.URL(URL_EXTRACTOR)
	extractor204.FrameType(FRAME_TYPE_EXTRACTOR)
	extractor204.Method(METHOD_EXTRACTOR)
	extractor204.Type(TYPE_EXTRACTOR)
	extractor204.RequestHeaders(REQUEST_HEADER_REFERER_EXTRACTOR)
	extractor204.Label(LABEL_EXTRACTOR)
	// Here you init new extractors

	return []Extractor{*extractor204}
}
