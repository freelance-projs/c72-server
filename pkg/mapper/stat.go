package mapper

import (
	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func ToTxLogDepartmentDto(m *model.LendingStat) dto.TxLogDepartment {
	if m == nil {
		return dto.TxLogDepartment{}
	}

	return dto.TxLogDepartment{
		ID:         m.ID,
		Department: m.Department,
		Lending:    int(m.Lending),
		Returned:   int(m.Returned),
		CreatedAt:  m.CreatedAt,
	}
}

func ToLogCompanyDto(m *model.TxLogCompany) dto.TxLogCompany {
	if m == nil {
		return dto.TxLogCompany{}
	}

	return dto.TxLogCompany{
		ID:        m.ID,
		Company:   m.Company,
		Washing:   int(m.Washing),
		Returned:  int(m.Returned),
		CreatedAt: m.CreatedAt,
	}
}
