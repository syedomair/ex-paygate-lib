package models

// APIRequestTest Public
type APIRequestTest struct {
	Serial          string `json:"serial"`
	RespVar         string `json:"resp_var"`
	IDVar           string `json:"id_var"`
	TitleVar        string `json:"title_var"`
	FindPath        string `json:"find_path"`
	Path            string `json:"path"`
	Body            string `json:"body"`
	RespStatus      string `json:"resp_status"`
	RespHTTPCode    int    `json:"resp_http_code"`
	RespErrCode     string `json:"resp_err_code"`
	RespTestVar     string `json:"resp_test_var"`
	RespTestVarCond string `json:"resp_test_var_cond"`
	RespTestValue   string `json:"resp_test_value"`
}

// APIResponseError Public
type APIResponseError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
