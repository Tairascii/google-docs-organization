package handler

type CreateOrgPayload struct {
	Title string `json:"title"`
}

type CreateOrgResponse struct {
	ID string `json:"id"`
}

type AddUserPayload struct {
	UserId string `json:"userId"`
	OrgId  string `json:"orgId"`
	Role   string `json:"role"`
}
