package entities

type Candidate struct {
	ID           string `json:"id" dynamodbav:"id"`
	CompositeKey string `json:"compositeKey" dynamodbav:"compositeKey"`
	DocumentID   string `json:"documentId" dynamodbav:"documentId"`
	DocumentName string `json:"documentName" dynamodbav:"documentName"`
	SheetID      string `json:"sheetId" dynamodbav:"sheetId"`
	SheetName    string `json:"sheetName" dynamodbav:"sheetName"`

	CreatedAt string  `json:"createdAt" dynamodbav:"createdAt"`
	CreatedBy string  `json:"createdBy" dynamodbav:"createdBy"`
	UpdatedAt *string `json:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty"`
	UpdatedBy *string `json:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty" dynamodbav:"deletedAt,omitempty"`
	DeletedBy *string `json:"deletedBy,omitempty" dynamodbav:"deletedBy,omitempty"`
	Deleted   bool    `json:"deleted" dynamodbav:"deleted"`
}
