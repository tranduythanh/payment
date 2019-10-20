package payment

import (
	"fmt"
	"net/url"

	"gopkg.in/go-playground/validator.v9"
)

// OnePayInternational ...
type OnePayInternational struct {
	PaymentGatewayHost string
	PaymentGatewayPath string
	Version            int
	Currency           string
	Command            string
	Locale             string

	Cfg *Config
}

// NewSandboxInternational ...
func NewSandboxInternational(returnURL string) *OnePayInternational {
	return &OnePayInternational{
		Version:  2,
		Currency: "VND",
		Command:  "pay",
		Locale:   "vn",

		Cfg: &Config{
			PaymentGatewayHost: "mtf.onepay.vn",
			PaymentGatewayPath: "vpcpay/vpcpay.op",
			QueryDRPath:        "vpcpay/Vpcdps.op",
			Merchant:           "TESTONEPAY",
			AccessCode:         "6BEB2546",
			SecureSecret:       "6D0870CDE5F24F34F3915FB0045120DB",
			ReturnURL:          returnURL,
		},
	}
}

// NewInternational ...
func NewInternational(cfg *Config) *OnePayInternational {
	return &OnePayInternational{
		Version:  2,
		Currency: "VND",
		Command:  "pay",
		Locale:   "vn",

		Cfg: cfg,
	}
}

// BuildCheckoutURL ...
func (op *OnePayInternational) BuildCheckoutURL(params *CheckoutParams) (string, error) {

	err := validator.New().Struct(params)
	if err != nil {
		return "", err
	}

	v := url.Values{}

	// Static params
	v.Add("vpc_Version", fmt.Sprintf("%d", op.Version))
	v.Add("vpc_Currency", op.Currency)
	v.Add("vpc_Command", op.Command)
	v.Add("vpc_AccessCode", op.Cfg.AccessCode)
	v.Add("vpc_Merchant", op.Cfg.Merchant)
	v.Add("vpc_Locale", op.Locale)
	v.Add("vpc_ReturnURL", op.Cfg.ReturnURL)

	// checkout params
	v.Add("vpc_MerchTxnRef", params.MerchTxnRef)
	v.Add("vpc_OrderInfo", params.OrderInfo)
	v.Add("vpc_Amount", fmt.Sprintf("%d00", params.Amount))
	v.Add("vpc_TicketNo", params.TicketNo)

	// Add SecureHash
	addSecureHash(&v, op.Cfg.SecureSecret)

	v.Add("Title", "Tran Duy Thanh")
	v.Add("AgainLink", "https://www.google.com.vn/")

	// Gen full url
	u := &url.URL{
		Scheme:   "https",
		Host:     op.Cfg.PaymentGatewayHost,
		Path:     op.Cfg.PaymentGatewayPath,
		RawQuery: v.Encode(),
	}

	return u.String(), nil
}

// HandleCallback ...
func (op *OnePayInternational) HandleCallback(v url.Values) (*InternationalResponse, error) {
	var resp = &InternationalResponse{}
	err := handleCallback(v, op.Cfg.SecureSecret, resp)
	if err != nil {
		return nil, err
	}

	resp.PostProcess()

	return resp, nil
}

// QueryDR ...Truy vấn trạng thái giao dịch (QueryDR API)
// - Chỉ gọi hàm này sau 15 phút giao dịch, Phương thức là redirect, kiểu GET
func (op *OnePayInternational) QueryDR(request *QueryDRAPIRequest) (res *QueryDRAPIResponse, err error) {
	return queryDR(op.Cfg, request)
}
