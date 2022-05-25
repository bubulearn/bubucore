package notifications

// AmoCRMAddLeadReq request to add new AmoCRM lead
type AmoCRMAddLeadReq struct {
	Lead    *AmoCRMLead    `json:"lead"`
	Contact *AmoCRMContact `json:"contact,omitempty"`
}

// AmoCRMItemWithID is a common object with ID
type AmoCRMItemWithID struct {
	ID int `json:"id"`
}

// AmoCRMCustomField is a custom field structure
// https://www.amocrm.ru/developers/content/crm_platform/custom-fields#cf-fill-examples
type AmoCRMCustomField struct {
	FieldID int         `json:"field_id"`
	Values  interface{} `json:"values"`
}

// AmoCRMLinks is a response links data
type AmoCRMLinks struct {
	Self *AmoCRMLink `json:"self"`
}

// AmoCRMLink is a response link data
type AmoCRMLink struct {
	Href string `json:"href"`
}

// region AmoCRMLead

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
	Embedded           *AmoCRMLeadEmbedded      `json:"_embedded,omitempty"`
}

// AmoCRMLeadCustomField is a custom field structure
// https://www.amocrm.ru/developers/content/crm_platform/custom-fields#cf-fill-examples
type AmoCRMLeadCustomField struct {
	AmoCRMCustomField
}

// AmoCRMLeadEmbedded is an embedded lead data
type AmoCRMLeadEmbedded struct {
	Tags      []*AmoCRMItemWithID          `json:"tags,omitempty"`
	Companies []*AmoCRMItemWithID          `json:"companies,omitempty"`
	Contacts  []*AmoCRMLeadEmbeddedContact `json:"contacts,omitempty"`
}

// AmoCRMLeadEmbeddedContact is an embedded lead contact
type AmoCRMLeadEmbeddedContact struct {
	AmoCRMItemWithID
	IsMain bool `json:"is_main"`
}

// endregion AmoCRMLead

// region AmoCRMContact

// AmoCRMContact is an info about contact
// https://www.amocrm.ru/developers/content/crm_platform/contacts-api#contacts-add
type AmoCRMContact struct {
	Name      string `json:"name,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`

	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	CreatedBy         int `json:"created_by,omitempty"`
	UpdatedBy         int `json:"updated_by,omitempty"`

	CreatedAt int `json:"created_at,omitempty"`
	UpdatedAt int `json:"updated_at,omitempty"`

	RequestID string `json:"request_id,omitempty"`

	CustomFieldsValues []*AmoCRMCustomField   `json:"custom_fields_values,omitempty"`
	Embedded           *AmoCRMContactEmbedded `json:"_embedded,omitempty"`
}

// AmoCRMContactEmbedded is an embedded contact data
type AmoCRMContactEmbedded struct {
	Tags []*AmoCRMItemWithID `json:"tags,omitempty"`
}

// AmoCRMContactsResp is an add contacts response
type AmoCRMContactsResp struct {
	Links    *AmoCRMLinks                `json:"_links"`
	Embedded *AmoCRMContactsEmbeddedResp `json:"_embedded"`
}

// AmoCRMContactsEmbeddedResp is an add contacts embedded response
type AmoCRMContactsEmbeddedResp struct {
	Contacts []*AmoCRMContactsEmbeddedContactResp `json:"contacts"`
}

// AmoCRMContactsEmbeddedContactResp is an embedded contact response
type AmoCRMContactsEmbeddedContactResp struct {
	ID        int          `json:"id"`
	RequestID string       `json:"request_id"`
	Links     *AmoCRMLinks `json:"_links"`
}

// endregion AmoCRMContact
