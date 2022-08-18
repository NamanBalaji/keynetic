package types

// key-val req res types
type GetSuccesResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type GetFailResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type DeleteSuccesResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
}

type DeleteFailResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PutRequest struct {
	Value string `json:"value"`
}

type PutSuccesResp struct {
	Replaced bool   `json:"replaced"`
	Message  string `json:"message,omitempty"`
}

type PutFailResp struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// views req res types
type GetViewSucces struct {
	Message string `json:"message,omitempty"`
	View    string `json:"view,omitempty"`
}

type DeleteViewSucces struct {
	Message string `json:"message,omitempty"`
}

type DeleteViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PutViewSucces struct {
	Message string `json:"message,omitempty"`
}

type PutViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// misc req res types
type GetStoreResponse struct {
	Store map[string]string `json:"store"`
}

type GetVectorClockResponse struct {
	VectorClock map[string]int `json:"vectorClock"`
}
