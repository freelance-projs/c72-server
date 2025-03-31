package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToSettingDto(m *model.Setting) dto.Setting {
	if m == nil {
		return dto.Setting{}
	}

	return dto.Setting{
		Key:   m.Key,
		Value: m.Value,
	}
}
