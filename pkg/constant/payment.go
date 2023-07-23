package constant

type PaymentStatus string

const (
	PAYMENT_STATUS_PENDING PaymentStatus = "PENDING"
	PAYMENT_STATUS_SUCCESS PaymentStatus = "SUCCESS"
	PAYMENT_STATUS_FAILED  PaymentStatus = "FAILED"
	PAYMENT_STATUS_EXPIRED PaymentStatus = "EXPIRED"
)

var mapPaymentStatus = map[PaymentStatus]string{
	PAYMENT_STATUS_PENDING: "PENDING",
	PAYMENT_STATUS_SUCCESS: "SUCCESS",
	PAYMENT_STATUS_FAILED:  "FAILED",
	PAYMENT_STATUS_EXPIRED: "EXPIRED",
}

func (s PaymentStatus) String() string {
	item, ok := mapPaymentStatus[s]
	if ok {
		return item
	}

	return "unknown"
}

type PaymentMethod string

const (
	PAYMENT_METHOD_VIRTUAL_ACCOUNT PaymentMethod = "Virtual Account"
	PAYMENT_METHOD_CREDIT_CARD     PaymentMethod = "Credit Card"
	PAYMENT_METHOD_BANK_TRANSFER   PaymentMethod = "Bank Transfer"
	PAYMENT_METHOD_E_WALLET        PaymentMethod = "e-Wallet"
	PAYMENT_METHOD_GOOGLE_WALLET   PaymentMethod = "Google Wallet"
)

var PaymentMethodString = map[string]PaymentMethod{
	"Virtual Account": PAYMENT_METHOD_VIRTUAL_ACCOUNT,
	"Credit Card":     PAYMENT_METHOD_CREDIT_CARD,
	"Bank Transfer":   PAYMENT_METHOD_BANK_TRANSFER,
	"e-Wallet":        PAYMENT_METHOD_E_WALLET,
	"Google Wallet":   PAYMENT_METHOD_GOOGLE_WALLET,
}

func GetPaymentMethod(key string) PaymentMethod {
	item, ok := PaymentMethodString[key]
	if ok {
		return item
	}

	return "unknown"
}
