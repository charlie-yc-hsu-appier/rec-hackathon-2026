package header

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type hmacHeaderTestSuite struct {
	suite.Suite
	mockClock *MockClock
}

func (ts *hmacHeaderTestSuite) SetupTest() {
	ts.mockClock = NewMockClock(gomock.NewController(ts.T()))
}

func TestHmacHeaderTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &hmacHeaderTestSuite{})
}

func (s *hmacHeaderTestSuite) TestGenerateHeaders() {
	s.mockClock.EXPECT().getDatetimeGMT().Return("250707T103117Z")

	header := &HmacHeader{
		SecretKey: "secret_key",
		AccessKey: "access_key",
		Clock:     s.mockClock,
	}

	result := header.GenerateHeaders(Params{
		RequestURL: "https://api-gateway.coupang.com/v2/providers/affiliate_open_api/apis/openapi/v1/products/reco?deviceId=FAKE-USER&imageSize=300x300&subId=KRpartner01",
	})

	wantedSignature := "CEA algorithm=HmacSHA256, access-key=access_key, signed-date=250707T103117Z, signature=04a3ea3f1e087fb00de58591123578c7080ba99efc24ebf337eb772ae7085023"
	require.Equal(s.T(), map[string]string{"Authorization": wantedSignature}, result)
}
