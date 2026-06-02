// 短信服务

package sms

import "context"

type Service interface {
  Send(ctx context.Context, tpl string, args []string, numbers ...string) error // number 不定长度方便传多个
}
