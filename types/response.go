package types

// var kk = "{'code': 428, 'message': 'The config needs to be initialized first.', 'ok': False, 'result': '/settings', 'time_taken': 1.6600999515503645e-05, 'title': 'Dester', 'description': 'Dester'}"

type DataResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Ok          bool   `json:"ok"`
	Result      any    `json:"result"`
	TimeTaken   int    `json:"time_taken"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type RootResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}
