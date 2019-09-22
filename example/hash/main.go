package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/k0kubun/pp"
)

var (
	secureSecret = "A3EFDFABA8653DF2342E8DAC29B51AF0"
)

func main() {

	sha := hash256(secureSecret, vpcParamsURL)

	sha = strings.ToUpper(sha)

	pp.Println(sha)
	pp.Println(hashStr)
}

func hash256(secret, urlQuery string) string {
	b, err := hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, b)
	h.Write([]byte(vpcParamsURL))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

var vpcParamsURL = `vpc_AccessCode=D67342C2&vpc_Amount=90000000&vpc_Command=pay&vpc_Currency=VND&vpc_Customer_Email=dev@naustud.io&vpc_Customer_Id=dev@naustud.io&vpc_Customer_Phone=0123456789&vpc_Locale=vn&vpc_MerchTxnRef=node-2019-09-21T02:28:15.302Z&vpc_Merchant=ONEPAY&vpc_OrderInfo=node-2019-09-21T02:28:15.302Z&vpc_ReturnURL=http://localhost:8080/payment/onepaydom/callback&vpc_SHIP_City=01&vpc_SHIP_Country=VN&vpc_SHIP_Provice=Hồ Chí Minh&vpc_SHIP_Street01=187 Dien Bien Phu, Da Kao Ward&vpc_TicketNo=::1&vpc_Version=2`
var hashStr = `CE24B16DDB3D1CA28B970370F7A8EDC82EAEC1E1801BFA710F00F41BE3705F3F`
