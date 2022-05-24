package notifications

// AmoCRMLead is an info about some client's request
// https://www.amocrm.ru/developers/content/crm_platform/leads-api#leads-add
type AmoCRMLead struct {
	Name string `json:"name,omitempty"`

	Price             int `json:"price,omitempty"`
	StatusID          int `json:"status_id,omitempty"`
	PipelineID        int `json:"pipeline_id,omitempty"`
	CreatedBy         int `json:"created_by,omitempty"`
	UpdatedBy         int `json:"updated_by,omitempty"`
	ClosedAt          int `json:"closed_at,omitempty"`
	CreatedAt         int `json:"created_at,omitempty"`
	UpdatedAt         int `json:"updated_at,omitempty"`
	LossReasonID      int `json:"loss_reason_id,omitempty"`
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`

	CustomFieldsValues []*AmoCRMLeadCustomField `json:"custom_fields_values,omitempty"`
}

// AmoCRMLeadCustomField is a custom field structure
// https://www.amocrm.ru/developers/content/crm_platform/custom-fields#cf-fill-examples
type AmoCRMLeadCustomField struct {
	FieldID int         `json:"field_id"`
	Values  interface{} `json:"values"`
}
