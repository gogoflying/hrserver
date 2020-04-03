package types

type ReqJoblist struct {
	City   string    `json:"city"`
	Page   int64     `json:"page"`
	Limit  int64     `json:"limit"`
	Search *KeyWords `json:"search_keywords"`
}

type KeyWords struct {
	SalaryHigh int64  `json:"salary_high"`
	SalaryLow  int64  `json:"salary_low"`
	Keys       string `json:"keys"`
}

type RspJoblist struct {
	Id          int64  `json:"job_id"`
	JobName     string `json:"job_name"`
	JobSalary   string `json:"job_salary"`
	JobCity     string `json:"job_city"`
	JobYears    string `json:"job_years"`
	JobEdu      string `json:"job_edu"`
	JobType     string `json:"job_type"`
	JobTime     string `json:"job_time"`
	CompanyName string `json:"company_name"`
	CompanyImg  string `json:"company_img_url"`
}
