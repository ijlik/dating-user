package repository

import (
	"database/sql"
	"time"

	"github.com/ijlik/dating-user/pkg/constant"
)

type OneTimePasswordLog struct {
	ID                  string                       `db:"id"`
	UserID              string                       `db:"user_id"`
	OneTimePasswordType constant.OneTimePasswordType `db:"onetime_password_type"`
	Code                string                       `db:"code"`
	Status              constant.OneTimeLogStatus    `db:"status"`
	CreatedAt           time.Time                    `db:"created_at"`
	UpdatedAt           sql.NullTime                 `db:"updated_at"`
	OTPLimit            int                          `db:"otp_limit"`
}

func (otl *OneTimePasswordLog) RowData() []interface{} {
	var data = []interface{}{
		otl.UserID,
		otl.OneTimePasswordType,
		otl.Code,
		otl.Status,
		otl.CreatedAt,
	}
	return data
}
