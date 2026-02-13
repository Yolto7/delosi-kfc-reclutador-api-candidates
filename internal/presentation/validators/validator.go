package validators

import (
	"github.com/Yolto7/api-candidates/internal/application/services/commands"
	"github.com/Yolto7/api-candidates/internal/application/services/queries"
	"github.com/Yolto7/api-candidates/pkg/infrastructure/validators"
)

// Validations
func GetByID(input queries.GetByIDServiceInput) error {
	return validators.ValidateSchema(&input)
}

func Create(input *commands.CreateServiceInput) error {
	return validators.ValidateSchema(input)
}

func Delete(input *commands.DeleteServiceInput) error {
	return validators.ValidateSchema(input)
}