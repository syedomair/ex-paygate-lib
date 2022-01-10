package models

// APIRequest Public
type APIRequest struct {
	RespVar  string `json:"resp_var"`
	IDVar    string `json:"id_var"`
	TitleVar string `json:"title_var"`
	FindPath string `json:"find_path"`
	Path     string `json:"path"`
	Body     string `json:"body"`
}
