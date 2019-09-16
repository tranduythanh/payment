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
		PaymentGatewayHost: "mtf.onepay.vn",
		PaymentGatewayPath: "vpcpay/vpcpay.op",
		Version:            2,
		Currency:           "VND",
		Command:            "pay",
		Locale:             "vn",

		Cfg: &Config{
			Merchant:     "TESTONEPAY",
			AccessCode:   "6BEB2546",
			SecureSecret: "6D0870CDE5F24F34F3915FB0045120DB",
			ReturnURL:    returnURL,
		},
	}
}

// NewInternational ...
func NewInternational(cfg *Config) *OnePayInternational {
	return &OnePayInternational{
		PaymentGatewayHost: "onepay.vn",
		PaymentGatewayPath: "vpcpay/vpcpay.op",
		Version:            2,
		Currency:           "VND",
		Command:            "pay",
		Locale:             "vn",

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

	// Gen full url
	u := &url.URL{
		Scheme:   "https",
		Host:     op.PaymentGatewayHost,
		Path:     op.PaymentGatewayPath,
		RawQuery: v.Encode(),
	}

	return u.String(), nil
}

// export function callbackOnePayInternational(req, res) {
// 	const query = req.query;

// 	return onepayIntl.verifyReturnUrl(query).then(results => {
// 		if (results) {
// 			res.locals.email = 'tu.nguyen@naustud.io';
// 			res.locals.orderId = results.orderId || '';
// 			res.locals.price = results.amount;
// 			res.locals.billingStreet = results.billingStreet;
// 			res.locals.billingCountry = Countries[results.billingCountry];
// 			res.locals.billingStateProvince = results.billingStateProvince;
// 			res.locals.billingCity = results.billingCity;
// 			res.locals.billingPostalCode = results.billingPostCode;

// 			res.locals.isSucceed = results.isSuccess;
// 			res.locals.message = results.message;
// 		} else {
// 			res.locals.isSucceed = false;
// 		}
// 	});
// }
