package types

// key-val req res types
type GetSuccesResp struct {
	Exists         bool   `json:"doesExist"`
	Message        string `json:"message,omitempty"`
	Value          string `json:"value,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
}

type GetFailResp struct {
	Exists         bool   `json:"doesExist"`
	Message        string `json:"message,omitempty"`
	Error          string `json:"error,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
}

type GetReq struct {
	CausalMetadata string `json:"causal-metadata"`
}

type DeleteSuccesResp struct {
	Exists         bool   `json:"doesExist"`
	Message        string `json:"message,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
}

type DeleteFailResp struct {
	Exists         bool   `json:"doesExist"`
	Message        string `json:"message,omitempty"`
	Error          string `json:"error,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
}

type DeleteReq struct {
	CausalMetadata string `json:"causal-metadata"`
}

type PutRequest struct {
	Value          string `json:"value"`
	CausalMetadata string `json:"causal-metadata"`
}

type PutSuccesResp struct {
	Replaced       bool   `json:"replaced"`
	Message        string `json:"message,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
}

type PutFailResp struct {
	Message        string `json:"message,omitempty"`
	Error          string `json:"error,omitempty"`
	CausalMetadata string `json:"causal-metadata"`
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
