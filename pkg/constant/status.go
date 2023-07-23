package constant

type UserStatus string

const (
	USER_STATUS_UNVERIFIED UserStatus = "UNVERIFIED"
	USER_STATUS_ACTIVE     UserStatus = "ACTIVE"
	USER_STATUS_DEACTIVE   UserStatus = "DEACTIVE"
)

var mapStatus = map[UserStatus]string{
	USER_STATUS_UNVERIFIED: "UNVERIFIED",
	USER_STATUS_ACTIVE:     "ACTIVE",
	USER_STATUS_DEACTIVE:   "DEACTIVE",
}

func (s UserStatus) String() string {
	item, ok := mapStatus[s]
	if ok {
		return item
	}

	return "unknown"
}
