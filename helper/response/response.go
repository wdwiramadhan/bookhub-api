package response

// ResponseError represent the reseponse error struct
type ResponseFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResponseSuccess represent the reseponse success struct
type ResponseSuccess struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
