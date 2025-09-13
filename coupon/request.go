package coupon

import (
	"encoding/json"
	"fmt"
)

type CreateCouponReq struct {
	Type    string        `json:"type"`
	Details CouponDetails `json:"details"`
}

func (r *CreateCouponReq) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type    CouponType      `json:"type"`
		Details json.RawMessage `json:"details"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Type = string(raw.Type)

	if raw.Type == "" || raw.Details == nil {
		return fmt.Errorf("invalid body, required field 'type' and 'details'")
	}

	switch raw.Type {
	case couponTypes[0]:
		var d CartWiseDetails
		if err := json.Unmarshal(raw.Details, &d); err != nil {
			return err
		}
		r.Details = d

	case couponTypes[1]:
		var d ProductWiseDetails
		if err := json.Unmarshal(raw.Details, &d); err != nil {
			return err
		}
		r.Details = d

	case couponTypes[2]:
		var d BxGyDetails
		if err := json.Unmarshal(raw.Details, &d); err != nil {
			return err
		}
		r.Details = d

	default:
		return fmt.Errorf("unsupported coupon type: %s", raw.Type)
	}

	return nil
}

func (r CreateCouponReq) Validate() error {
	if r.Type == "" {
		return fmt.Errorf("type is required field")
	}
	if r.Details == nil {
		return fmt.Errorf("details is required field")
	}
	return r.Details.ValidateCoupon()
}
