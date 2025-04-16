package dto

type GetDepartmentStatRequest struct {
	Department string `json:"department" uri:"department" validate:"required"`
}

type ListDepartmentStatRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type DepartmentStat struct {
	Department string `json:"department" header:"Phòng ban"`
	Exported   int    `json:"exported" header:"Đang mượn"`
	Returned   int    `json:"returned" header:"Trả"`
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
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

type CompanyStat struct {
	Company  string `json:"company" header:"Công ty"`
	Exported int    `json:"exported" header:"Đang giặt"`
	Returned int    `json:"returned" header:"Trả"`
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
	TagName         string `json:"tag_name" header:"Tên vật phẩm"`
	Lending         int    `json:"lending" header:"Đang mượn"`
	LendingReturned int    `json:"lending_returned" header:"Khoa trả"`
	Washing         int    `json:"washing" header:"Đang giặt"`
	WashingReturned int    `json:"washing_returned" header:"Công ty trả"`
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

type ListTagStatRequest struct {
	From *int64 `json:"from" form:"from" validate:"required,omitnil"`
	To   *int64 `json:"to" form:"to" validate:"required,omitnil"`
}

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
