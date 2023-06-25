package extractor

import "tds/shared/models"

type ExtractorTypes int

const (
	DocumentId ExtractorTypes = iota
	DocumentLifecycle
	FrameId
	FrameType
	Initiator
	Method
	ParentFrameId
	RequestId
	TabId
	TimeStamp
	Type
	URL
	Success
	RequestHeaders
	Label
)

type DocumentIdExtractor func(string) []int
type DocumentLifecycleExtractor func(string) []int
type FrameIdExtractor func(int) []int
type FrameTypeExtractor func(string) []int
type InitiatorExtractor func(string) []int
type MethodExtractor func(string) []int
type ParentFrameIdExtractor func(int) []int
type RequestIdExtractor func(string) []int
type TabIdExtractor func(string) []int
type TimeStampExtractor func(string) []int
type TypeExtractor func(string) []int
type URLExtractor func(string) []int
type SuccessExtractor func(string) []int
type RequestHeadersExtractor func([]map[string]string) []int
type LabelExtractor func(bool) []int

type Extractor struct {
	sequence                   []ExtractorTypes
	documentIdExtractor        DocumentIdExtractor
	documentLifecycleExtractor DocumentLifecycleExtractor
	frameIdExtractor           FrameIdExtractor
	frameTypeExtractor         FrameTypeExtractor
	initiatorExtractor         InitiatorExtractor
	methodExtractor            MethodExtractor
	parentFrameIdExtractor     ParentFrameIdExtractor
	requestIdExtractor         RequestIdExtractor
	tabIdExtractor             TabIdExtractor
	timeStampExtractor         TimeStampExtractor
	typeExtractor              TypeExtractor
	urlExtractor               URLExtractor
	successExtractor           SuccessExtractor
	requestHeadersExtractor    RequestHeadersExtractor
	labelExtractor             LabelExtractor
}

func NewExtractor() *Extractor {
	return &Extractor{
		sequence: make([]ExtractorTypes, 0),
	}
}

func (e *Extractor) DocumentId(extractor DocumentIdExtractor) {
	e.documentIdExtractor = extractor
	e.sequence = append(e.sequence, DocumentId)
}

func (e *Extractor) DocumentLifecycle(extractor DocumentLifecycleExtractor) {
	e.documentLifecycleExtractor = extractor
	e.sequence = append(e.sequence, DocumentLifecycle)
}

func (e *Extractor) FrameId(extractor FrameIdExtractor) {
	e.frameIdExtractor = extractor
	e.sequence = append(e.sequence, FrameId)
}

func (e *Extractor) FrameType(extractor FrameTypeExtractor) {
	e.frameTypeExtractor = extractor
	e.sequence = append(e.sequence, FrameType)
}

func (e *Extractor) Initiator(extractor InitiatorExtractor) {
	e.initiatorExtractor = extractor
	e.sequence = append(e.sequence, Initiator)
}

func (e *Extractor) Method(extractor MethodExtractor) {
	e.methodExtractor = extractor
	e.sequence = append(e.sequence, Method)
}

func (e *Extractor) ParentFrameId(extractor ParentFrameIdExtractor) {
	e.parentFrameIdExtractor = extractor
	e.sequence = append(e.sequence, ParentFrameId)
}

func (e *Extractor) RequestId(extractor RequestIdExtractor) {
	e.requestIdExtractor = extractor
	e.sequence = append(e.sequence, RequestId)
}

func (e *Extractor) TabId(extractor TabIdExtractor) {
	e.tabIdExtractor = extractor
	e.sequence = append(e.sequence, TabId)
}

func (e *Extractor) TimeStamp(extractor TimeStampExtractor) {
	e.timeStampExtractor = extractor
	e.sequence = append(e.sequence, TimeStamp)
}

func (e *Extractor) TypeExtractor(extractor TypeExtractor) {
	e.typeExtractor = extractor
	e.sequence = append(e.sequence, Type)
}

func (e *Extractor) URL(extractor URLExtractor) {
	e.urlExtractor = extractor
	e.sequence = append(e.sequence, URL)
}

func (e *Extractor) Success(extractor SuccessExtractor) {
	e.successExtractor = extractor
	e.sequence = append(e.sequence, Success)
}

func (e *Extractor) RequestHeaders(extractor RequestHeadersExtractor) {
	e.requestHeadersExtractor = extractor
	e.sequence = append(e.sequence, RequestHeaders)
}

func (e *Extractor) Label(extractor LabelExtractor) {
	e.labelExtractor = extractor
	e.sequence = append(e.sequence, Label)
}

func (e *Extractor) Encode(requestData models.RequestData) []int {
	return make([]int, 0)
}
