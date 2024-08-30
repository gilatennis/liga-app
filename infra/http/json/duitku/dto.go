package duitku

type CallbackPayload struct {
	MerchantCode     string `form:"merchantCode"`
	Amount           int    `form:"amount"`
	MerchantOrderID  string `form:"merchantOrderId"`
	ProductDetail    string `form:"productDetail"`
	AdditionalParam  string `form:"additionalParam"`
	ResultCode       string `form:"resultCode"`
	Signature        string `form:"signature"`
	PaymentCode      string `form:"paymentCode"`
	MerchantUserID   string `form:"merchantUserId"`
	Reference        string `form:"reference"`
	SpUserHash       string `form:"spUserHash"`
	SettlementDate   string `form:"settlementDate"`
	PublisherOrderID string `form:"publisherOrderId"`
	SourceAccount    string `form:"sourceAccount"`
	// Other fields are omitted
}
