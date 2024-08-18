package v2

type V2EmailOutboundEmailInput struct {
	ConfigurationSetName                      *string                `json:"ConfigurationSetName"`
	Content                                   Content                `json:"Content"`
	Destination                               Destination            `json:"Destination"`
	EmailTags                                 []EmailTag             `json:"EmailTags"`
	FeedbackForwardingEmailAddress            *string                `json:"FeedbackForwardingEmailAddress"`
	FeedbackForwardingEmailAddressIdentityArn *string                `json:"FeedbackForwardingEmailAddressIdentityArn"`
	FromEmailAddress                          *string                `json:"FromEmailAddress"`
	FromEmailAddressIdentityArn               *string                `json:"FromEmailAddressIdentityArn"`
	ListManagementOptions                     *ListManagementOptions `json:"ListManagementOptions"`
	ReplyToAddresses                          []string               `json:"ReplyToAddresses"`
}

type Content struct {
	Raw      *Raw      `json:"Raw"`
	Simple   *Simple   `json:"Simple"`
	Template *Template `json:"Template"`
}

type Raw struct {
	Data string `json:"Data"`
}

type Simple struct {
	Body    Body    `json:"Body"`
	Subject Subject `json:"Subject"`
}

type Body struct {
	Html *Html `json:"Html"`
	Text *Text `json:"Text"`
}

type Html struct {
	CharSet *string `json:"CharSet"`
	Data    string  `json:"Data"`
}

type Text struct {
	CharSet *string `json:"CharSet"`
	Data    string  `json:"Data"`
}

type Subject struct {
	CharSet *string `json:"CharSet"`
	Data    string  `json:"Data"`
}

type Template struct {
	TemplateArn  *string `json:"TemplateArn"`
	TemplateData *string `json:"TemplateData"`
	TemplateName *string `json:"TemplateName"`
}

type Destination struct {
	BccAddresses []string `json:"BccAddresses"`
	CcAddresses  []string `json:"CcAddresses"`
	ToAddresses  []string `json:"ToAddresses"`
}

type EmailTag struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type ListManagementOptions struct {
	ContactListName *string `json:"ContactListName"`
	TopicName       *string `json:"TopicName"`
}
