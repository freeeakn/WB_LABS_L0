package order

import "strings"

func (o *Order) Validate() error {
	if o == nil {
		return ErrInvalidOrder
	}
	if strings.TrimSpace(o.OrderUID) == "" {
		return ErrMissingField("order_uid")
	}
	if strings.TrimSpace(o.TrackNumber) == "" {
		return ErrMissingField("track_number")
	}
	return nil
}
