package codes

const (
	Success               = 200
	Created               = 201	// 创建成功
	Bad_Request           = 400 // 请求错误
	Unauthorized          = 401 // 没有登录
	Not_Found             = 404 // not found
	Forbidden             = 403 // 没有权限
	Method_Not_Allowed    = 405 // 方法不对 (POST,PUT,GET)
	Not_Acceptable        = 406 // 不能通过
	Internal_Server_Error = 500 // 服务错误

	Params_Err = 430 // 参数错误
	Capt_Err   = 431 // 验证码错误
)

func CodeString(code int32) string {
	s := map[int32]string{
		Success:               "OK",
		Bad_Request:           "Bad_Request",
		Unauthorized:          "Unauthorized",
		Not_Found:             "Not_Found",
		Forbidden:             "Forbidden",
		Method_Not_Allowed:    "Method_Not_Allowed",
		Not_Acceptable:        "Not_Acceptable",
		Internal_Server_Error: "Server_Error",
		Params_Err:            "Params_Error",
		Capt_Err:              "Captcha_Error",
	}[code]
	return s
}
