package payment

import (
	"net/url"
	"testing"

	"github.com/k0kubun/pp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOnePay(t *testing.T) {
	Convey("OnePay", t, func() {
		Convey("validateHash", func() {
			sample := `https://a98f94fa.ngrok.io/payment/callback/international?vpc_OrderInfo=1569179952150041000&vpc_3DSECI=02&vpc_AVS_Street01=Tran+Quang+Khai&vpc_Merchant=TESTONEPAY&vpc_Card=MC&vpc_AcqResponseCode=00&AgainLink=https%3A%2F%2Fwww.google.com.vn%2F&vpc_AVS_Country=VN&vpc_AuthorizeId=023209&vpc_3DSenrolled=Y&vpc_ReceiptNo=926519189281&vpc_TransactionNo=011EN9&vpc_AVS_StateProv=Hanoi&vpc_Locale=vn&vpc_TxnResponseCode=0&vpc_VerToken=jHyn%2B7YFi1EUAREAAAAvNUe6Hv8%3D&vpc_Amount=10000000&vpc_BatchNo=2.0190923E7&vpc_Version=2&vpc_AVSResultCode=Y&vpc_VerStatus=M&vpc_Command=pay&vpc_Message=Transaction+Approved&Title=Tran+Duy+Thanh&vpc_3DSstatus=Y&vpc_SecureHash=B233867E915FFC67A2EA71E79E4C7FED44B589FD7AD71DD57ADA5642247EBA50&vpc_CardNum=531358xxxxxxx430&vpc_AVS_PostCode=1234&vpc_CSCResultCode=MATCH&vpc_MerchTxnRef=1569179952150041000&vpc_VerType=3DS&vpc_3DSXID=Io6FdLRj91zPA%2BYZum4s19HrNkU%3D&vpc_AVS_City=North&vpc_CommercialCardIndicator=1`

			u, err := url.Parse(sample)
			So(err, ShouldBeNil)

			v := u.Query()
			pp.Println(len(v))
			ok, err := validateSecureHash(&v, "6D0870CDE5F24F34F3915FB0045120DB")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})
	})
}
