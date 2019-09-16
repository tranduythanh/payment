package payment

import (
	"fmt"
	"net/url"

	"gopkg.in/go-playground/validator.v9"
)

// OnePayDomestic ...
type OnePayDomestic struct {
	PaymentGatewayHost string
	PaymentGatewayPath string
	Version            int
	Currency           string
	Command            string
	Locale             string

	Cfg *Config
}

// NewSandboxDomestic ...
func NewSandboxDomestic(returnURL string) *OnePayDomestic {
	return &OnePayDomestic{
		PaymentGatewayHost: "mtf.onepay.vn",
		PaymentGatewayPath: "onecomm-pay/vpc.op",
		Version:            2,
		Currency:           "VND",
		Command:            "pay",
		Locale:             "vn",

		Cfg: &Config{
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
		PaymentGatewayHost: "onepay.vn",
		PaymentGatewayPath: "onecomm-pay/vpc.op",
		Version:            2,
		Currency:           "VND",
		Command:            "pay",
		Locale:             "vn",

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
		Host:     op.PaymentGatewayHost,
		Path:     op.PaymentGatewayPath,
		RawQuery: v.Encode(),
	}

	return u.String(), nil
}

// export function callbackOnePayDomestic(req, res) {
// 	const query = req.query;

// 	return onepayDom.verifyReturnUrl(query).then(results => {
// 		if (results) {
// 			res.locals.email = 'tu.nguyen@naustud.io';
// 			res.locals.orderId = results.orderId || '';
// 			res.locals.price = results.amount;

// 			res.locals.isSucceed = results.isSuccess;
// 			res.locals.message = results.message;
// 		} else {
// 			res.locals.isSucceed = false;
// 		}
// 	});
// }
