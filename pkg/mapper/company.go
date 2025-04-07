package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToCompanyDto(m *model.Company) dto.Company {
	if m == nil {
		return dto.Company{}
	}

	return dto.Company{
		Name: m.Name,
	}
}
