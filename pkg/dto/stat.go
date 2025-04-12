package dto

type GetDepartmentStatRequest struct {
	Department string `json:"department" uri:"department" validate:"required"`
}

type ListDepartmentStatRequest struct {
}

type DepartmentStat struct {
	Department string `json:"department"`
	Exported   int    `json:"exported"`
	Returned   int    `json:"returned"`
}

type DepartmentDetailStat struct {
	Department string        `json:"department"`
	Exported   int           `json:"exported"`
	Returned   int           `json:"returned"`
	Trackings  []TagTracking `json:"tracking"`
}

// company
type GetCompanyStatRequest struct {
	Company string `json:"company" uri:"company" validate:"required"`
}

type ListCompanyStatRequest struct {
}

type CompanyStat struct {
	Company  string `json:"company"`
	Exported int    `json:"exported"`
	Returned int    `json:"returned"`
}

type CompanyDetailStat struct {
	Company   string        `json:"company"`
	Exported  int           `json:"exported"`
	Returned  int           `json:"returned"`
	Trackings []TagTracking `json:"tracking"`
}

// tag
type GetTagStatRequest struct {
	TagName string `json:"tag_name" uri:"tag_name" validate:"required"`
}

type TagStat struct {
	TagName         string `json:"tag_name"`
	Lending         int    `json:"lending"`
	LendingReturned int    `json:"lending_returned"`
	Washing         int    `json:"washing"`
	WashingReturned int    `json:"washing_returned"`
}

type TagStatDetail struct {
	TagName         string               `json:"tag_name"`
	Lending         int                  `json:"lending"`
	LendingReturned int                  `json:"lending_returned"`
	Washing         int                  `json:"washing"`
	WashingReturned int                  `json:"washing_returned"`
	Departments     []DepartmentTracking `json:"departments"`
	Companies       []CompanyTracking    `json:"companies"`
}

type ListTagStatRequest struct{}

type DepartmentTracking struct {
	Name     string `json:"name"`
	Exported int    `json:"exported"`
	Returned int    `json:"returned"`
}

type CompanyTracking struct {
	Name     string `json:"name"`
	Exported int    `json:"exported"`
	Returned int    `json:"returned"`
}
