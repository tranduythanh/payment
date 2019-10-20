package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/go-playground/validator.v9"
)

// Defines ...
const (
	VPCSecureHashKey = "vpc_SecureHash"
	VPCPrefix        = "vpc_"
)

// Config ...
type Config struct {
	PaymentGatewayHost string `validate:"required" yaml:"payment_gateway_host" json:"payment_gateway_host"`
	PaymentGatewayPath string `validate:"required" yaml:"payment_gateway_path" json:"payment_gateway_path"`
	QueryDRPath        string `validate:"required" yaml:"query_dr_path" json:"query_dr_path"`
	Merchant           string `validate:"required" yaml:"merchant" json:"merchant"`
	AccessCode         string `validate:"required" yaml:"access_code" json:"access_code"`
	ReturnURL          string `validate:"required,max=128" yaml:"return_url" json:"return_url"`
	SecureSecret       string `validate:"required" yaml:"secure_secret" json:"secure_secret"`
}

// CheckoutParams ...
type CheckoutParams struct {
	Amount      int64  `validate:"required,lte=9999999999"`
	OrderInfo   string `validate:"required,max=34"`
	MerchTxnRef string `validate:"required,max=40"`
	TicketNo    string `validate:"required,max=15"`
	Title       string `validate:"required,max=64"`
	AgainLink   string `validate:"required,max=64"`
}

// How to gen secure hash
// - all url params wit prefix 'vpc_', sorted by name asc
// - SECURE_SECRET provided by ONEPAY
func addSecureHash(v *url.Values, secureSecret string) error {
	data, err := genStringForHash(v)
	if err != nil {
		return err
	}

	sha, err := genHash(data, secureSecret)
	if err != nil {
		return err
	}

	v.Set(VPCSecureHashKey, sha)

	return nil
}

func validateSecureHash(v *url.Values, secureSecret string) (bool, error) {
	data, err := genStringForHash(v)
	if err != nil {
		return false, err
	}

	receivedSecureHash := v.Get(VPCSecureHashKey)

	secureHash, err := genHash(data, secureSecret)
	if err != nil {
		return false, err
	}

	return receivedSecureHash == secureHash, nil
}

func genStringForHash(v *url.Values) (string, error) {
	vpcParams := &url.Values{}
	for key := range *v {
		if !strings.HasPrefix(key, VPCPrefix) {
			continue
		}
		if key == VPCSecureHashKey {
			continue
		}
		vpcParams.Add(key, v.Get(key))
	}
	return url.QueryUnescape(vpcParams.Encode())
}

func genHash(urlQuery, secret string) (string, error) {
	unescapedQuery, err := url.QueryUnescape(urlQuery)
	if err != nil {
		return "", err
	}

	hexByteSecret, err := hex.DecodeString(secret)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, hexByteSecret)
	h.Write([]byte(unescapedQuery))
	sha := hex.EncodeToString(h.Sum(nil))

	return strings.ToUpper(sha), nil
}

func handleCallback(v url.Values, secureSecret string, resp interface{}) error {

	ok, err := validateSecureHash(&v, secureSecret)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Invalid secure_hash")
	}

	var decoder = schema.NewDecoder()

	decoder.IgnoreUnknownKeys(true)

	err = decoder.Decode(resp, map[string][]string(v))
	if err != nil {
		return err
	}

	return nil
}

// DomesticResponse ...
type DomesticResponse struct {
	VPCAdditionData    string `json:"vpc_AdditionData" query:"vpc_AdditionData" schema:"vpc_AdditionData"`
	VPCCommand         string `json:"vpc_Command" query:"vpc_Command" schema:"vpc_Command"`
	VPCCurrencyCode    string `json:"vpc_CurrencyCode" query:"vpc_CurrencyCode" schema:"vpc_CurrencyCode"`
	VPCLocale          string `json:"vpc_Locale" query:"vpc_Locale" schema:"vpc_Locale"`
	VPCMerchTxnRef     string `json:"vpc_MerchTxnRef" query:"vpc_MerchTxnRef" schema:"vpc_MerchTxnRef"`
	VPCTransactionNo   string `json:"vpc_TransactionNo" query:"vpc_TransactionNo" schema:"vpc_TransactionNo"`
	VPCVersion         string `json:"vpc_Version" query:"vpc_Version" schema:"vpc_Version"`
	VPCSecureHash      string `json:"vpc_SecureHash" query:"vpc_SecureHash" schema:"vpc_SecureHash"`
	VPCAcqResponseCode string `json:"vpc_AcqResponseCode" query:"vpc_AcqResponseCode" schema:"vpc_AcqResponseCode"`
	VPCAuthorizeID     string `json:"vpc_AuthorizeId" query:"vpc_AuthorizeId" schema:"vpc_AuthorizeId"`
	VPCCard            string `json:"vpc_Card" query:"vpc_Card" schema:"vpc_Card"`
	VPCCardNum         string `json:"vpc_CardNum" query:"vpc_CardNum" schema:"vpc_CardNum"`
	VPCCommercialCard  string `json:"vpc_CommercialCard" query:"vpc_CommercialCard" schema:"vpc_CommercialCard"`

	VPC3DSECI      string `json:"vpc_3DSECI" query:"vpc_3DSECI" schema:"vpc_3DSECI"`
	VPC3Dsenrolled string `json:"vpc_3Dsenrolled" query:"vpc_3Dsenrolled" schema:"vpc_3Dsenrolled"`
	VPC3Dsstatus   string `json:"vpc_3Dsstatus" query:"vpc_3Dsstatus" schema:"vpc_3Dsstatus"`

	VPCMessage         string `json:"vpc_Message" query:"vpc_Message" schema:"vpc_Message"`
	VPCMerchant        string `json:"vpc_Merchant" query:"vpc_Merchant" schema:"vpc_Merchant"`
	VPCAmount          int64  `json:"vpc_Amount" query:"vpc_Amount" schema:"vpc_Amount"`
	VPCOrderInfo       string `json:"vpc_OrderInfo" query:"vpc_OrderInfo" schema:"vpc_OrderInfo"`
	VPCTxnResponseCode string `json:"vpc_TxnResponseCode" query:"vpc_TxnResponseCode" schema:"vpc_TxnResponseCode"`

	AgainLink string `json:"AgainLink" query:"AgainLink" schema:"AgainLink"`
	Title     string `json:"Title" query:"Title" schema:"Title"`

	TxnResponseMessage ErrorMessageLocale `json:"txnResponseCode" query:"txnResponseCode" schema:"txnResponseCode"`
}

// InternationalResponse ...
type InternationalResponse struct {
	VPCAdditionData  string `json:"vpc_AdditionData" query:"vpc_AdditionData" schema:"vpc_AdditionData"`
	VPCCommand       string `json:"vpc_Command" query:"vpc_Command" schema:"vpc_Command"`
	VPCCurrencyCode  string `json:"vpc_CurrencyCode" query:"vpc_CurrencyCode" schema:"vpc_CurrencyCode"`
	VPCLocale        string `json:"vpc_Locale" query:"vpc_Locale" schema:"vpc_Locale"`
	VPCMerchTxnRef   string `json:"vpc_MerchTxnRef" query:"vpc_MerchTxnRef" schema:"vpc_MerchTxnRef"`
	VPCTransactionNo string `json:"vpc_TransactionNo" query:"vpc_TransactionNo" schema:"vpc_TransactionNo"`
	VPCVersion       string `json:"vpc_Version" query:"vpc_Version" schema:"vpc_Version"`
	VPCSecureHash    string `json:"vpc_SecureHash" query:"vpc_SecureHash" schema:"vpc_SecureHash"`

	VPCCard           string `json:"vpc_Card" query:"vpc_Card" schema:"vpc_Card"`
	VPCCardNum        string `json:"vpc_CardNum" query:"vpc_CardNum" schema:"vpc_CardNum"`
	VPCCommercialCard string `json:"vpc_CommercialCard" query:"vpc_CommercialCard" schema:"vpc_CommercialCard"`

	VPC3DSECI      string `json:"vpc_3DSECI" query:"vpc_3DSECI" schema:"vpc_3DSECI"`
	VPC3DSXID      string `json:"vpc_3DSXID" query:"vpc_3DSXID" schema:"vpc_3DSXID"`
	VPC3DSenrolled string `json:"vpc_3DSenrolled" query:"vpc_3DSenrolled" schema:"vpc_3DSenrolled"`
	VPC3DSstatus   string `json:"vpc_3DSstatus" query:"vpc_3DSstatus" schema:"vpc_3DSstatus"`

	VPCAVSStreet01   string `json:"vpc_AVS_Street01" query:"vpc_AVS_Street01" schema:"vpc_AVS_Street01"`
	VPCAVSCountry    string `json:"vpc_AVS_Country" query:"vpc_AVS_Country" schema:"vpc_AVS_Country"`
	VPCAVSStateProv  string `json:"vpc_AVS_StateProv" query:"vpc_AVS_StateProv" schema:"vpc_AVS_StateProv"`
	VPCAVSCity       string `json:"vpc_AVS_City" query:"vpc_AVS_City" schema:"vpc_AVS_City"`
	VPCAVSResultCode string `json:"vpc_AVSResultCode" query:"vpc_AVSResultCode" schema:"vpc_AVSResultCode"`
	VPCAVSPostCode   string `json:"vpc_AVS_PostCode" query:"vpc_AVS_PostCode" schema:"vpc_AVS_PostCode"`

	VPCAcqResponseCode         string `json:"vpc_AcqResponseCode" query:"vpc_AcqResponseCode" schema:"vpc_AcqResponseCode"`
	VPCAuthorizeID             string `json:"vpc_AuthorizeId" query:"vpc_AuthorizeId" schema:"vpc_AuthorizeId"`
	VPCRiskOverallResult       string `json:"vpc_RiskOverallResult" query:"vpc_RiskOverallResult" schema:"vpc_RiskOverallResult"`
	VPCReceiptNo               string `json:"vpc_ReceiptNo" query:"vpc_ReceiptNo" schema:"vpc_ReceiptNo"`
	VPCBatchNo                 string `json:"vpc_BatchNo" query:"vpc_BatchNo" schema:"vpc_BatchNo"`
	VPCCSCResultCode           string `json:"vpc_CSCResultCode" query:"vpc_CSCResultCode" schema:"vpc_CSCResultCode"`
	VPCCommercialCardIndicator string `json:"vpc_CommercialCardIndicator" query:"vpc_CommercialCardIndicator" schema:"vpc_CommercialCardIndicator"`

	VPCVerType          string `json:"vpc_VerType" query:"vpc_VerType" schema:"vpc_VerType"`
	VPCVerToken         string `json:"vpc_VerToken" query:"vpc_VerToken" schema:"vpc_VerToken"`
	VPCVerStatus        string `json:"vpc_VerStatus" query:"vpc_VerStatus" schema:"vpc_VerStatus"`
	VPCVerSecurityLevel string `json:"vpc_VerSecurityLevel" query:"vpc_VerSecurityLevel" schema:"vpc_VerSecurityLevel"`

	VPCMessage         string `json:"vpc_Message" query:"vpc_Message" schema:"vpc_Message"`
	VPCMerchant        string `json:"vpc_Merchant" query:"vpc_Merchant" schema:"vpc_Merchant"`
	VPCAmount          int64  `json:"vpc_Amount" query:"vpc_Amount" schema:"vpc_Amount"`
	VPCOrderInfo       string `json:"vpc_OrderInfo" query:"vpc_OrderInfo" schema:"vpc_OrderInfo"`
	VPCTxnResponseCode string `json:"vpc_TxnResponseCode" query:"vpc_TxnResponseCode" schema:"vpc_TxnResponseCode"`

	AgainLink string `json:"AgainLink" query:"AgainLink" schema:"AgainLink"`
	Title     string `json:"Title" query:"Title" schema:"Title"`

	TxnResponseMessage ErrorMessageLocale `json:"txnResponseCode" query:"txnResponseCode" schema:"txnResponseCode"`
}

// PostProcess ...
func (r *DomesticResponse) PostProcess() {
	r.VPCAmount = r.VPCAmount / 100
	r.TxnResponseMessage = ErrorMap[r.VPCTxnResponseCode]
}

// PostProcess ...
func (r *InternationalResponse) PostProcess() {
	r.VPCAmount = r.VPCAmount / 100
	r.TxnResponseMessage = ErrorMap[r.VPCTxnResponseCode]
}

// QueryDRAPIRequest ...
type QueryDRAPIRequest struct {
	VPCCommand       string `json:"vpc_Command" query:"vpc_Command" schema:"vpc_Command"`
	VPCVersion       string `json:"vpc_Version" query:"vpc_Version" schema:"vpc_Version"`
	VPCMerchTxnRef   string `json:"vpc_MerchTxnRef" query:"vpc_MerchTxnRef" schema:"vpc_MerchTxnRef"`
	VPCMerchant      string `json:"vpc_Merchant" query:"vpc_Merchant" schema:"vpc_Merchant"`
	VPCAccessCode    string `json:"vpc_AccessCode" query:"vpc_AccessCode" schema:"vpc_AccessCode"`
	VPCUser          string `json:"vpc_User" query:"vpc_User" schema:"vpc_User"`
	VPCPassword      string `json:"vpc_Password" query:"vpc_Password" schema:"vpc_Password"`
	VPCSecureHashKey string `json:"vpc_SecureHash" query:"vpc_SecureHash" schema:"vpc_SecureHash"`
}

// QueryDRAPIResponse ...
type QueryDRAPIResponse struct {
	VPCDRExists        string `json:"vpc_DRExists" query:"vpc_DRExists" schema:"vpc_DRExists"`
	VPCTxnResponseCode string `json:"vpc_TxnResponseCodes" query:"vpc_TxnResponseCodes" schema:"vpc_TxnResponseCodes"`
}

func queryDR(cfg *Config, request *QueryDRAPIRequest) (res *QueryDRAPIResponse, err error) {
	if cfg == nil {
		return nil, fmt.Errorf("Config is nil")
	}

	err = validator.New().Struct(request)
	if err != nil {
		return nil, err
	}

	v := url.Values{}

	v.Add("vpc_Command", request.VPCCommand)
	v.Add("vpc_Version", request.VPCVersion)
	v.Add("vpc_MerchTxnRef", request.VPCMerchTxnRef)
	v.Add("vpc_Merchant", request.VPCMerchant)
	v.Add("vpc_AccessCode", request.VPCAccessCode)
	v.Add("vpc_User", request.VPCUser)
	v.Add("vpc_Password", request.VPCPassword)

	addSecureHash(&v, cfg.SecureSecret)

	u := &url.URL{
		Scheme:   "https",
		Host:     cfg.PaymentGatewayHost,
		Path:     cfg.QueryDRPath,
		RawQuery: v.Encode(),
	}

	_, body, errs := gorequest.New().Get(u.String()).End()
	if len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	return res, err
}
