package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

// Config ...
type Config struct {
	Merchant     string `validate:"required"`
	AccessCode   string `validate:"required"`
	ReturnURL    string `validate:"required,max=128"`
	SecureSecret string `validate:"required"`
}

// CheckoutParams ...
type CheckoutParams struct {
	Amount      int64  `validate:"required,lte=9999999999"`
	OrderInfo   string `validate:"required,max=34"`
	MerchTxnRef string `validate:"required,max=40"`
	TicketNo    string `validate:"required,max=15"`
}

// Qui tắc tạo chữ ký:
// - Đầu vào là phần chuỗi tham số từ dấu “?” sau Payment URL hoặc ReturnURL
// - Chuỗi chỉ chứa các tham số có tiền tố vpc_
// - Chuỗi tham số sắp xếp tham số theo thứ tự anphabe
// - Key dùng để Hash là chuỗi hexa do OnePAY cấp cho mỗi Merchant ID (gọi là
//   hash code hoặc SECURE_SECRET)
func addSecureHash(v *url.Values, secureSecret string) {
	urlQuery := v.Encode()
	h := hmac.New(sha256.New, []byte(secureSecret))
	h.Write([]byte(urlQuery))
	sha := hex.EncodeToString(h.Sum(nil))
	sha = strings.ToUpper(sha)
	v.Add("vpc_SecureHash", sha)
}

// DomesticResponse ...
type DomesticResponse struct {
	VPCMessage         string `json:"vpc_Message" query:"vpc_Message"`
	VPCMerchant        string `json:"vpc_Merchant" query:"vpc_Merchant"`
	VPCAmount          int    `json:"vpc_Amount" query:"vpc_Amount"`
	VPCOrderInfo       string `json:"vpc_OrderInfo" query:"vpc_OrderInfo"`
	VPCTxnResponseCode int    `json:"vpc_TxnResponseCode" query:"vpc_TxnResponseCode"`
}

// Norm ...
func (r *DomesticResponse) Norm() {
	r.VPCAmount = r.VPCAmount / 100
}
