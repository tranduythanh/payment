package payment

import (
	"fmt"
	"net/url"

	"gopkg.in/go-playground/validator.v9"
)

// OnePayDomestic ...
type OnePayDomestic struct {
	Version  int
	Currency string
	Command  string
	Locale   string

	Cfg *Config
}

// NewSandboxDomestic ...
func NewSandboxDomestic(returnURL string) *OnePayDomestic {
	return &OnePayDomestic{
		Version:  2,
		Currency: "VND",
		Command:  "pay",
		Locale:   "vn",

		Cfg: &Config{
			PaymentGatewayHost: "mtf.onepay.vn",
			PaymentGatewayPath: "onecomm-pay/vpc.op",

			Merchant:     "ONEPAY",
			AccessCode:   "D67342C2",
			SecureSecret: "A3EFDFABA8653DF2342E8DAC29B51AF0",
			ReturnURL:    returnURL,
		},
	}
}

// NewDomestic ...
func NewDomestic(cfg *Config) *OnePayDomestic {
	return &OnePayDomestic{
		Version:  2,
		Currency: "VND",
		Command:  "pay",
		Locale:   "vn",

		Cfg: cfg,
	}
}

// BuildCheckoutURL ...
func (op *OnePayDomestic) BuildCheckoutURL(params *CheckoutParams) (string, error) {

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
	v.Add("AgainLink", "localhost:8080/payment/checkout")

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
func (op *OnePayDomestic) HandleCallback(v url.Values) (*DomesticResponse, error) {
	var resp = &DomesticResponse{}
	err := handleCallback(v, op.Cfg.SecureSecret, resp)
	if err != nil {
		return nil, err
	}

	resp.PostProcess()

	return resp, nil
}
