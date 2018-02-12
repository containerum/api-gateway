package mail

//Types related to sending emails

type SimpleSendRequest struct {
	Template  string                 `json:"template" binding:"required"`
	UserID    string                 `json:"user_id" binding:"required,uuid4"`
	Variables map[string]interface{} `json:"variables"`
}

type SimpleSendResponse struct {
	UserID string `json:"user_id"`
}

type Recipient struct {
	ID        string                 `json:"id" binding:"required,uuid4"`
	Name      string                 `json:"name" binding:"required"`
	Email     string                 `json:"email" binding:"required,email"`
	Variables map[string]interface{} `json:"variables"`
}

type SendRequest struct {
	Delay   int                `json:"delay" binding:"omitempty,min=0"` // in minutes
	Message struct {
		CommonVariables map[string]string `json:"common_variables"`
		Recipients      []Recipient       `json:"recipient_data" binding:"required"`
	} `json:"message" binding:"required"`
}

type SendStatus struct {
	RecipientID  string `json:"recipient_id"`
	TemplateName string `json:"template_name"`
	Status       string `json:"status"`
}

type SendResponse struct {
	Statuses []SendStatus `json:"email_list"`
}