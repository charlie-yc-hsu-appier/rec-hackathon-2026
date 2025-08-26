package unmarshaler

import (
	"encoding/json"
	"fmt"
)

type CoupangPartner struct{}

func (s *CoupangPartner) UnmarshalResponse(body []byte) (*[]PartnerResp, error) {
	var resp []PartnerResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
	}
	return &resp, nil
}
