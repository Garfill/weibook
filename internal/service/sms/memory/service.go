package memory

import (
  "context"
  "fmt"
)

type Service struct {
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
  fmt.Println(numbers)
  return nil
}

func NewMemorySms() *Service {
  return &Service{}
}
