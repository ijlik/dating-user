package repository

import "github.com/ijlik/dating-user/pkg/constant"

type Payment struct {
	UserID        string                 `db:"user_id"`
	Amount        float32                `db:"amount"`
	Identifier    string                 `db:"identifier"`
	PaymentMethod constant.PaymentMethod `db:"payment_method"`
	PaymentData   string                 `db:"payment_data"`
	Status        constant.PaymentStatus `db:"status"`
}

func (p *Payment) RowData() []interface{} {
	var data = []interface{}{
		p.UserID,
		p.Amount,
		p.Identifier,
		p.PaymentMethod,
		p.PaymentData,
		p.Status,
	}
	return data
}
