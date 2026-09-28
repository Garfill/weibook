package res

// 通用返回结果包装
// Code 业务码

type Result struct {
  Code int    `json:"code"`
  Msg  string `json:"msg"`
  Data any    `json:"data"`
}
