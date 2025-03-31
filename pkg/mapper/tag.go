package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToTagDto(m *model.Tag) dto.Tag {
	if m == nil {
		return dto.Tag{}
	}

	return dto.Tag{
		ID:   m.ID,
		Name: m.Name.String,
	}
}
