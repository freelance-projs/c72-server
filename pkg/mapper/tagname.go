package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToTagNameDto(m *model.TagName) dto.TagName {
	if m == nil {
		return dto.TagName{}
	}

	return dto.TagName{
		Name: m.Name,
	}
}
