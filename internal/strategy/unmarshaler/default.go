package unmarshaler

import (
	"encoding/json"
	"fmt"
)

type Default struct{}

func (s *Default) UnmarshalResponse(body []byte) (*[]CoupangPartnerResp, error) {
	var resp []CoupangPartnerResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
	}
	return &resp, nil
}
