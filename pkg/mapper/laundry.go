package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToLaundryDto(m *model.Laundry) dto.Laundry {
	if m == nil {
		return dto.Laundry{}
	}

	return dto.Laundry{
		ID:         m.ID,
		Name:       m.Name,
		NumWashing: m.NumWashing,
		CreatedAt:  m.CreatedAt,
	}
}

func ToLaundryTagDto(m *model.LaundryTag) dto.LaundryTag {
	if m == nil {
		return dto.LaundryTag{}
	}

	return dto.LaundryTag{
		LaundryID: m.LaundryID,
		TagID:     m.TagID,
		Status:    m.Status.String(),
	}
}
