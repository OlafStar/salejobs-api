package types

type DecodedToken struct {
	Username string `json:"username"`
	Exp int64 `json:"exp"`
}

type User struct {
	Username string
	Password string
}

//Advertisment types

type AdvertismentCounterResponse struct {
	Total int64 `json:"total"`
}

type Company struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}

type SalaryRange struct {
	EmploymentType string `json:"employmentType"`
	MinSalary      int64  `json:"minSalary"`
	MaxSalary      int64  `json:"maxSalary"`
	Currency       string `json:"currency"`
}

type Location struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type ContactDetails struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ApplyTypeEnum struct {
	Email string `json:"email"`
	Url   string `json:"url"`
}

type CreateAdvertisementBody struct {
	Company       Company        `json:"company"`
	Title         string         `json:"title"`
	Experience    string         `json:"experience"`
	Skill         string         `json:"skill"`
	Salary        []SalaryRange  `json:"salary"`
	Description   string         `json:"description"`
	Location      Location       `json:"location"`
	OperatingMode string         `json:"operatingMode"`
	TypeOfWork    string         `json:"typeOfWork"`
	ApplyType     ApplyTypeEnum  `json:"applyType"`
	Consent       bool           `json:"consent"`
	Contact       ContactDetails `json:"contact"`
}

type CreateAdvertisementResponse struct {
	Id            string         `json:"id"`
	Company       Company        `json:"company"`
	Title         string         `json:"title"`
	Experience    string         `json:"experience"`
	Skill         string         `json:"skill"`
	Salary        []SalaryRange  `json:"salary"`
	Description   string         `json:"description"`
	Location      Location       `json:"location"`
	OperatingMode string         `json:"operatingMode"`
	TypeOfWork    string         `json:"typeOfWork"`
	ApplyType     ApplyTypeEnum  `json:"applyType"`
	Consent       bool           `json:"consent"`
	Contact       ContactDetails `json:"contact"`
}

type GetAdvertismentBody struct {
	Page int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type GetAdvertismentResponse struct {
	CurrentPage int64 `json:"currentPage"`
	Total int64 `json:"total"`
	Last int64 `json:"last"`
	Advertisments []CreateAdvertisementResponse `json:"advertisments"`
}