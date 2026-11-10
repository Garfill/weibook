package ratelimit

import (
  "context"
  "testing"
  "weibook/internal/rateLimit"
  limitermocks "weibook/internal/rateLimit/mock"
  "weibook/internal/service/sms"
  smsmocks "weibook/internal/service/sms/mock"

  "github.com/stretchr/testify/assert"
  "go.uber.org/mock/gomock"
)

func TestRateLimitSmsService_Send(t *testing.T) {
  tcs := []struct {
    name string

    mock func(c *gomock.Controller) (sms.Service, rateLimit.Limiter)

    wantErr error
  }{
    {
      name: "限流",
      mock: func(c *gomock.Controller) (sms.Service, rateLimit.Limiter) {
        smsS := smsmocks.NewMockService(c)
        limit := limitermocks.NewMockLimiter(c)
        limit.EXPECT().Limited(gomock.Any(), gomock.Any()).Return(true, nil)
        return smsS, limit
      },
      wantErr: ErrSMSRateLimit,
    },
  }

  for _, tt := range tcs {
    t.Run(tt.name, func(t *testing.T) {
      ctrl := gomock.NewController(t)
      defer ctrl.Finish()

      sms, l := tt.mock(ctrl)
      svc := NewRateLimitSmsService(sms, l)
      err := svc.Send(context.Background(), "abc", []string{"123123"}, "123456")
      assert.Equal(t, tt.wantErr, err)
    })
  }

}
