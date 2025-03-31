package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToDepartmentDto(m *model.Department) dto.Department {
	if m == nil {
		return dto.Department{}
	}

	return dto.Department{
		Name: m.Name,
	}
}
