package v1

type SendEmailInput struct {
	Version              string
	ConfigurationSetName string
	ToAddresses          string
	CcAddresses          string
	BccAddresses         string
	HtmlData             string
	HtmlCharset          string
	TextData             string
	TextCharset          string
	SubjectData          string
	SubjectCharset       string
	ReplyToAddresses     string
	ReturnPath           string
	ReturnPathArn        string
	Source               string
	SourceArn            string
	Tags                 string
	Destination          string
	FromArn              string
}
