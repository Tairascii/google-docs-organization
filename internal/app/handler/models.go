package handler

type CreateOrgPayload struct {
	Title string `json:"title"`
}

type CreateOrgResponse struct {
	ID string `json:"id"`
}

type AddUserPayload struct {
	UserEmail string `json:"userEmail"`
	OrgId     string `json:"orgId"`
	Role      string `json:"role"`
}
