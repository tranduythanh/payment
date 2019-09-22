package payment

// ErrorMessageLocale ...
type ErrorMessageLocale struct {
	VN string
	EN string
}

// ErrorMap ...
var ErrorMap = map[int]ErrorMessageLocale{
	0: ErrorMessageLocale{
		VN: "Giao dịch thành công",
		EN: "Approved",
	},
	1: ErrorMessageLocale{
		VN: "Giao dịch không thành công. Ngân hàng phát hành thẻ từ chối cấp phép cho giao dịch. Vui lòng liên hệ ngân hàng theo số điện thoại sau mặt thẻ để biết chính xác nguyên nhân Ngân hàng từ chối.",
		EN: "The transaction is unsuccessful. This transaction has been declined by issuer bank. Please contact your bank for further clarification.",
	},
	3: ErrorMessageLocale{
		VN: "Mã đơn vị không tồn tại",
		EN: "Merchant not exist",
	},
	4: ErrorMessageLocale{
		VN: "Không đúng access code",
		EN: "Invalid access code",
	},
	5: ErrorMessageLocale{
		VN: "Số tiền không hợp lệ",
		EN: "Invalid amount",
	},
	6: ErrorMessageLocale{
		VN: "Mã tiền tệ không tồn tại",
		EN: "Invalid currency code",
	},
	7: ErrorMessageLocale{
		VN: "Lỗi không xác định",
		EN: "Unspecified Failure ",
	},
	8: ErrorMessageLocale{
		VN: "Số thẻ không đúng",
		EN: "Invalid card Number",
	},
	9: ErrorMessageLocale{
		VN: "Tên chủ thẻ không đúng",
		EN: "Invalid card name",
	},
	10: ErrorMessageLocale{
		VN: "Thẻ hết hạn/Thẻ bị khóa",
		EN: "Expired Card",
	},
	11: ErrorMessageLocale{
		VN: "Thẻ chưa đăng ký sử dụng dịch vụ",
		EN: "Card Not Registed Service(internet banking)",
	},
	12: ErrorMessageLocale{
		VN: "Ngày phát hành/Hết hạn không đúng",
		EN: "Invalid card date",
	},
	13: ErrorMessageLocale{
		VN: "Vượt quá hạn mức thanh toán",
		EN: "Exist Amount",
	},
	21: ErrorMessageLocale{
		VN: "Số tiền không đủ để thanh toán",
		EN: "Insufficient fund",
	},
	22: ErrorMessageLocale{
		VN: "Thông tin tài khoản không đúng",
		EN: "Invalid Account",
	},
	23: ErrorMessageLocale{
		VN: "Tài khoản bị khóa",
		EN: "Account Locked",
	},
	24: ErrorMessageLocale{
		VN: "Thông tin thẻ không đúng",
		EN: "Invalid Card Info",
	},
	25: ErrorMessageLocale{
		VN: "OTP không đúng",
		EN: "Invalid OTP",
	},
	253: ErrorMessageLocale{
		VN: "Quá thời gian thanh toán",
		EN: "Transaction timeout",
	},
	99: ErrorMessageLocale{
		VN: "Người sử dụng hủy giao dịch",
		EN: "User cancel",
	},
}
