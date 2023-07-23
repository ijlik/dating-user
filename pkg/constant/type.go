package constant

type OneTimePasswordType string

const (
	OTP_TYPE_EMAIL      OneTimePasswordType = "EMAIL"
	OTP_TYPE_SMS        OneTimePasswordType = "SMS"
	OTP_TYPE_PHONE_CALL OneTimePasswordType = "PHONE_CALL"
	OTP_TYPE_WHATSAPP   OneTimePasswordType = "WHATSAPP"
)

var mapOtpType = map[OneTimePasswordType]string{
	OTP_TYPE_EMAIL:      "EMAIL",
	OTP_TYPE_SMS:        "SMS",
	OTP_TYPE_PHONE_CALL: "PHONE_CALL",
	OTP_TYPE_WHATSAPP:   "WHATSAPP",
}

func (o OneTimePasswordType) String() string {
	item, ok := mapOtpType[o]
	if ok {
		return item
	}

	return "unknown"
}

type OneTimeLogStatus string

const (
	OTP_STATUS_UNUSED  OneTimeLogStatus = "UNUSED"
	OTP_STATUS_EXPIRED OneTimeLogStatus = "EXPIRED"
	OTP_STATUS_USED    OneTimeLogStatus = "USED"
)

var mapOtpStatus = map[OneTimeLogStatus]string{
	OTP_STATUS_UNUSED:  "UNUSED",
	OTP_STATUS_EXPIRED: "EXPIRED",
	OTP_STATUS_USED:    "USED",
}

func (o OneTimeLogStatus) String() string {
	item, ok := mapOtpStatus[o]
	if ok {
		return item
	}

	return "unknown"
}
