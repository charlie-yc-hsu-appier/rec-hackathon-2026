package header

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type replaceHeaderTestSuite struct {
	suite.Suite
	mockClock *MockClock
}

func (ts *replaceHeaderTestSuite) SetupTest() {
	ts.mockClock = NewMockClock(gomock.NewController(ts.T()))
}

func TestReplaceHeaderTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &replaceHeaderTestSuite{})
}

func (s *replaceHeaderTestSuite) TestGenerateHeaders() {
	s.mockClock.EXPECT().getDatetimeGMT().Return("250707T103117Z")

	header := &ReplaceHeader{
		SecretKey: "secret_key",
		AccessKey: "access_key",
		Clock:     s.mockClock,
	}

	result := header.GenerateHeaders(Params{
		RequestURL: "https://api-gateway.coupang.com/v2/providers/affiliate_open_api/apis/openapi/v2/products/reco",
		HTTPMethod: "POST",
	})

	wantedSignature := "CEA algorithm=HmacSHA256, access-key=access_key, signed-date=250707T103117Z, signature=faf13b58f6cc013892a036b778465bdcb85326d418c11398139a0b80ade01624"
	require.Equal(s.T(), map[string]string{"Authorization": wantedSignature}, result)
}
