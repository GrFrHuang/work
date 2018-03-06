package controllers

// unused
type AuthViewController struct {
	BaseController
}

func (c *AuthViewController) Login() {
	// todo 登录页面模板
	c.Data["json"] = "todo 这是登录页面"
	c.ServeJSON()
}
