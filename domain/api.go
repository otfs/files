package domain

//-------------------------------------------------------------
// Api接口层
//-------------------------------------------------------------

const (
	ResultCodeOk            = "0"
	ResultCodeBadRequest    = "400"
	ResultCodeInternalError = "500"
)

type Result struct {
	Code string      `json:"code"` // error code
	Msg  string      `json:"msg"`  // error message
	Data interface{} `json:"data"` // data
}

func NewResultOk(data interface{}) *Result {
	return &Result{
		Code: ResultCodeOk,
		Msg:  "ok",
		Data: data,
	}
}

func NewResultFail(code string, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
