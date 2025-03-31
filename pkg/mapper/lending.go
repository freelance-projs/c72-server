package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToLendingDTO(m *model.Lending) dto.Lending {
	if m == nil {
		return dto.Lending{}
	}

	return dto.Lending{
		ID:         m.ID,
		Department: m.Department,
		NumLending: m.NumLending,
		CreatedAt:  m.CreatedAt,
	}
}

func ToLendingTagDTO(m *model.LendingTag) dto.LendingTag {
	if m == nil {
		return dto.LendingTag{}
	}

	return dto.LendingTag{
		LendingID: m.LendingID,
		TagID:     m.TagID,
		Status:    m.Status.String(),
	}
}
