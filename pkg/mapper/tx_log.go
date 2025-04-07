package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToTxLogDept(m *model.TxLogDepartment) dto.TxLogDepartment {
	if m == nil {
		return dto.TxLogDepartment{}
	}

	return dto.TxLogDepartment{
		ID:          m.ID,
		Department:  m.Overview.Actor,
		NumLending:  int(m.Overview.TotalTags - m.Overview.Returned),
		NumReturned: int(m.Overview.Returned),
		CreatedAt:   m.CreatedAt,
	}
}

func ToTxLogCompany(m *model.TxLogCompany) dto.TxLogCompany {
	if m == nil {
		return dto.TxLogCompany{}
	}

	return dto.TxLogCompany{
		ID:          m.ID,
		Company:     m.Overview.Actor,
		NumWashing:  int(m.Overview.TotalTags - m.Overview.Returned),
		NumReturned: int(m.Overview.Returned),
		CreatedAt:   m.CreatedAt,
	}
}
