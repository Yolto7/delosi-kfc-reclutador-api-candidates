package entities

type Candidate struct {
	ID                                string `json:"id" dynamodbav:"id"`
	CompositeKey                      string `json:"compositeKey" dynamodbav:"compositeKey"`
	SheetID                           string `json:"sheetId" dynamodbav:"sheetId"`
	ColumnPostulantResponseId         string `json:"columnPostulantResponseId" dynamodbav:"columnPostulantResponseId"`
	ColumnPostulantDateTimeResponseId string `json:"columnPostulantDateTimeResponseId" dynamodbav:"columnPostulantDateTimeResponseId"`
	ColumnPostulantConfirmedId        string `json:"columnPostulantConfirmedId" dynamodbav:"columnPostulantConfirmedId"`

	CreatedAt string  `json:"createdAt" dynamodbav:"createdAt"`
	CreatedBy string  `json:"createdBy" dynamodbav:"createdBy"`
	UpdatedAt *string `json:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty"`
	UpdatedBy *string `json:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty" dynamodbav:"deletedAt,omitempty"`
	DeletedBy *string `json:"deletedBy,omitempty" dynamodbav:"deletedBy,omitempty"`
	Deleted   bool    `json:"deleted" dynamodbav:"deleted"`
}
