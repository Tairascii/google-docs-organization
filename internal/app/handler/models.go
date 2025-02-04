package handler

type CreateOrgPayload struct {
	Title string `json:"title"`
}

type CreateOrgResponse struct {
	ID string `json:"id"`
}
