package mail

import "time"

//Types related to storing templates

type TemplateStorageValue struct {
	Data      string    `json:"data"`
	Subject   string    `json:"template_subject"`
	CreatedAt time.Time `json:"created_at"` // UTC
}

type TemplateCreateRequest struct {
	Name    string `json:"template_name" binding:"required"`
	Version string `json:"template_version" binding:"required"`
	Data    string `json:"template_data" binding:"required,base64"`
	Subject string `json:"template_subject" binding:"required"`
}

type TemplateCreateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplateUpdateRequest struct {
	Data    string `json:"template_data" binding:"omitempty,base64"`
	Subject string `json:"template_subject" binding:"omitempty"`
}

type TemplateUpdateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplateDeleteResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplatesDeleteResponse struct {
	Name string `json:"template_name"`
}

type TemplatesListEntry struct {
	Name     string   `json:"template_name"`
	Versions []string `json:"template_versions"`
}

type TemplatesListResponse struct {
	Templates []TemplatesListEntry `json:"templates"`
}

type MessageListQuery struct {
	Page    int `form:"page" binding:"required,gt=0"`
	PerPage int `form:"per_page" binding:"required,gt=0"`
}
