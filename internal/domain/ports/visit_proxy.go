package ports

import "context"

type VisitCreateRequest struct {
	Type                string `json:"type" validate:"type"`
	BuildingID          string `json:"buildingId" validate:"uuid4,notblank"`
	UserID              string `json:"userId" validate:"uuid4,notblank"`
	VisitorFullName     string `json:"visitorFullName" dynamodbav:"notblank"`
	VisitDate           string `json:"visitDate" validate:"notblank"`
	Relationship 				string `json:"relationship" validate:"notblank"`
}

type VisitProxy interface {
	Create(ctx context.Context, input *VisitCreateRequest) error
}