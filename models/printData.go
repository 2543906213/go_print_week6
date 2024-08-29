package models

type PrintData struct {
	Printdata   []string      `json:"print_data"`
	Template    Template      `json:"template"`
	TemplateDtl []TemplateDtl `json:"template_dtl"`
	BoxData     []string      `json:"box_data"`
	Box         string        `json:"box"`
	PrintNum    int           `json:"print_num"`
}

type Template struct {
	TemplateId   string `json:"template_id"`
	TemplateName string `json:"template_name"`
	LabelWidth   string `json:"label_width"`
	LabelHeight  string `json:"label_height"`
	LabelMark    string `json:"label_mark"`
	FontsSize    string `json:"fonts_size"`
	LabelStyle   string `json:"label_style"`
	LabelType    string `json:"label_type"`
	Logo         string `json:"logo"`
	LogoLeft     string `json:"logo_left"`
	LogoUp       string `json:"logo_up"`
	LogoRight    string `json:"logo_right"`
	LogoDown     string `json:"logo_down"`
	QrcodeDown   string `json:"qrcode_down"`
	QrcodeRight  string `json:"qrcode_right"`
	QrcodeUp     string `json:"qrcode_up"`
	QrcodeLeft   string `json:"qrcode_left"`
	ShowLogo     string `json:"show_logo"`
	ShowQrcode   string `json:"show_qrcode"`
	HaveChapter  int    `json:"have_chapter"`
}

type TemplateDtl struct {
	TemplateId  string `json:"template_id"`
	Field       string `json:"field"`
	FieldDesc   string `json:"field_desc"`
	DescType    string `json:"desc_type"`
	Type        string `json:"type"`
	SameLine    string `json:"same_line"`
	MoveLeft    string `json:"move_left"`
	MoveUp      string `json:"move_up"`
	FieldLeft   string `json:"field_left"`
	FieldUp     string `json:"field_up"`
	FieldRight  string `json:"field_right"`
	FieldDown   string `json:"field_down"`
	QrcodeLeft  int
	QrcodeUp    int
	QrcodeRight int
	QrcodeDown  int
}

type RowData struct {
	Mailno         string
	SortingCode    string
	Receiver       string
	ReceiveAddress string
	Sender         string
	SenderAddress  string
	FromCompany    string
	GoodsName      string
	GoodsRemark    string
	Remark         string
	Qrcode         string
	Field          string
}
