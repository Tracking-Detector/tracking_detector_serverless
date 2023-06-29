package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestData struct {
	Id                primitive.ObjectID  `json:"_id,omitempty"`
	DocumentId        string              `json:"documentId"  validate:"required"`
	DocumentLifecycle string              `json:"documentLifecycle"  validate:"required"`
	FrameId           int                 `json:"frameId"  validate:"required"`
	FrameType         string              `json:"frameType"  validate:"required"`
	Initiator         string              `json:"initiator"  validate:"required"`
	Method            string              `json:"method"  validate:"required"`
	ParentFrameId     int                 `json:"parentFrameId"  validate:"required"`
	RequestId         string              `json:"requestId"  validate:"required"`
	TabId             int                 `json:"tabId"  validate:"required"`
	TimeStamp         string              `json:"timeStamp"  validate:"required"`
	Type              string              `json:"type"  validate:"required"`
	URL               string              `json:"url"  validate:"required"`
	RequestHeaders    []map[string]string `json:"requestHeaders"  validate:"required"`
	Success           bool                `json:"success"  validate:"required"`
	Label             bool                `json:"label"  validate:"required"`
}
