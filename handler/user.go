package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const pwd_salt = "*#890"

// SignupHandler: 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	} else {
		r.ParseForm()
		username := r.Form.Get("username")
		passwd := r.Form.Get("password")
		if len(username) < 3 || len(passwd) < 5 {
			w.Write([]byte("invalid parameter"))
			return
		}
		enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
		suc := dblayer.UserSignup(username, enc_passwd)
		if suc {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	}
}

// SignupHandler: 处理用户注册请求
//func SignupHandler(c *gin.Context) {
//	c.Redirect(http.StatusFound, "/static/view/signup.html")
//}

// DoSignupHandler: 处理注册的POST请求
//func DoSignupHandler(c *gin.Context) {
//	username := c.Request.FormValue("username")
//	passwd := c.Request.FormValue("password")
//	if len(username) < 3 || len(passwd) < 5 {
//		c.JSON(http.StatusOK, gin.H{
//			"msg":  "invalid parameter",
//			"code": -1,
//		})
//		return
//	}
//	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
//	suc := dblayer.UserSignup(username, enc_passwd)
//	if suc {
//		c.JSON(http.StatusOK, gin.H{
//			"msg":  "SUCCESS",
//			"code": 0,
//		})
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"msg":  "FAILED",
//			"code": -2,
//		})
//	}
//}

// SignInHandler: 登陆接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	} else {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		encPasswd := util.Sha1([]byte(password + pwd_salt))

		//1.校验用户名及密码
		pwdChecked := dblayer.UserSignin(username, encPasswd)
		if !pwdChecked {
			w.Write([]byte("FAILED"))
			return
		}

		//2.生成访问凭证(token)
		token := GenToken(username)
		upRes := dblayer.UpdateToken(username, token)
		if !upRes {
			w.Write([]byte("FAILED"))
			return
		}
		fmt.Println(r.Host)
		//3.登陆成功后重定向到首页
		//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
		resp := util.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: struct {
				Location string
				Username string
				Token    string
			}{
				Location: "http://" + r.Host + "/static/view/home.html",
				Username: username,
				Token:    token,
			},
		}
		w.Write(resp.JSONBytes())
	}
}

// UserInfoHandler： 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//拥有拦截器以后不需要校验
	//token := r.Form.Get("token")
	//
	//// 2. 验证token是否有效
	//isValidToken := IsTokenValid(username, token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且相应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())

}

// GenToken: 获取用户token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid: token是否有效
func IsTokenValid(usename string, token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenTime, _ := strconv.ParseInt(token[len(token)-8:], 16, 64)
	nowTime, _ := strconv.ParseInt(ts, 16, 64)
	if nowTime-tokenTime >= 10*365*12*60*60 {
		return false
	}
	if !dblayer.IsTokenEqual(usename, token) {
		return false
	}
	return true
}
