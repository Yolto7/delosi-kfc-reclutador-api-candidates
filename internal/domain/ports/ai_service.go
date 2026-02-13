package ports

import "context"

type AnalyzeCandidateRequest struct {
	Candidate 			string
	CurrentDateTime   string
	TimeZone          string
}

type DetectedVisit struct {
	VisitorFullName string `json:"visitorFullName"`
	VisitDate       string `json:"visitDate"`
	Relationship 		string `json:"relationship"`
}

type AnalyzeCandidateResponse struct {
	Visits []DetectedVisit
}

type AIService interface {
	AnalyzeCandidate(ctx context.Context, req AnalyzeCandidateRequest) (*AnalyzeCandidateResponse, error)
}

