package protocol

const (
	MethodCheckURLs = "check_urls"
)

type ProfileRequest struct {
	URLs []string `json:"urls"`
}

type ProfileResponse struct {
	URLCodes map[string]int `json:"url_codes,omitempty"`
	Error    string         `json:"error,omitempty"`
}
