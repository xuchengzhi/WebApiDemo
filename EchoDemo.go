package main

import "github.com/labstack/echo"
import "net/http"

// import "fmt"
// import "strconv"
// import "time"
import "encoding/base64"
import "encoding/hex"
import "crypto/sha256"
import "crypto/hmac"

// import "github.com/garyburd/redigo/redis"

// import "os"
// import "reflect"
// import "encoding/json"

func HMAC_SHA256(src, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return hex.EncodeToString(m.Sum(nil))
}

// base编码
func BASE64EncodeStr(src string) string {
	return string(base64.StdEncoding.EncodeToString([]byte(src)))
}

// base解码
func BASE64DecodeStr(src string) string {
	a, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(a)
}

func XorEncodeStr(msg, key string) string {
	ml := len(msg)
	kl := len(key)
	pwd := ""
	for i := 0; i < ml; i++ {
		pwd += (string((key[i%kl]) ^ (msg[i])))
	}

	return pwd
}

func XorDecodeStr(msg, key string) string {
	ml := len(msg)
	kl := len(key)
	pwd := ""
	for i := 0; i < ml; i++ {
		pwd += (string(((msg[i]) ^ key[i%kl])))
	}
	return pwd
}

type ApiJson struct {
	Status int         `json:"result"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"info"`
}

func ApiResource(status int, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: objects, Msg: msg}
	return
}

type StsJson struct {
	Times string `json:"t"`
	Sys   string `json:"sys"`
	Types string `json:"type"`
	Token string `json:"token"`
}

var is_xor = true

// func StsToken_v2(c echo.Context) error {
// 	data := make(map[string]string)
// 	sts_info := c.FormValue("sts_info")
// 	key := "xczceshi"

// 	if len(sts_info) == 0 {
// 		data["result"] = "5"
// 		data["msg"] = "error"
// 		data["info"] = "sts_info not in params"
// 	} else {

// 		msg := XorDecodeStr(sts_info, key)
// 		now := (time.Now().UnixNano() / 1e6)
// 		var jsons StsJson
// 		if err := json.Unmarshal([]byte(msg), &jsons); err == nil {
// 			t, errs := strconv.Atoi(jsons.Times)
// 			if errs != nil {
// 			}
// 			tt := int64(t)
// 			if now-tt > 1000000 {
// 				data["result"] = "1"
// 				data["msg"] = "error"
// 				data["info"] = "time out"
// 			} else {

// 				tmp := HMAC_SHA256(jsons.Times+jsons.Types+jsons.Sys, "ceshi1234")
// 				tokens := BASE64EncodeStr(tmp)

// 				if tokens == jsons.Token {
// 					toks := Redis_()
// 					data["result"] = "0"
// 					data["msg"] = "success"
// 					data["info"] = toks
// 				} else {
// 					data["result"] = "2"
// 					data["msg"] = "error"
// 					data["info"] = "json 校验失败"
// 				}
// 			}
// 		} else {
// 			data["result"] = "4"
// 			data["msg"] = "error"
// 			data["info"] = "Xor Msg Error"
// 		}
// 	}
// 	var XorMs string
// 	// fmt.Println("data= %d", data)
// 	jsonStr, err := json.Marshal(data)
// 	if err != nil {
// 	}
// 	if is_xor {
// 		XorMs = XorEncodeStr(string(jsonStr), key)
// 	} else {
// 		XorMs = string(jsonStr)
// 	}

// 	return c.String(http.StatusOK, string(XorMs))
// }

// func ApiTest(c echo.Context) error {
// 	// name := c.FormValue("name")
// 	times := time.Now()
// 	return c.JSON(http.StatusOK, ApiResource(0, times, "success"))
// }

// func StsToken(c echo.Context) error {
// 	apimsg := ""
// 	var apicode int
// 	var data string
// 	sts_info := c.FormValue("sts_info")
// 	key := "xczceshi"
// 	msg := XorDecodeStr(sts_info, key)

// 	var jsons StsJson
// 	if err := json.Unmarshal([]byte(msg), &jsons); err == nil {
// 		tmp := HMAC_SHA256(jsons.Times+jsons.Types+jsons.Sys, "ceshi1234")
// 		tokens := BASE64EncodeStr(tmp)

// 		if tokens == jsons.Token {
// 			toks := Redis_()
// 			apimsg = toks //"验证通过"
// 			apicode = 0
// 			data = "success"
// 		}

// 	} else {
// 		// ctx.StatusCode(iris.StatusBadRequest)
// 		apimsg = "json 校验失败"
// 		apicode = 2
// 		data = "error"
// 	}
// 	// ctx.StatusCode(iris.StatusOK)
// 	return c.JSON(http.StatusOK, ApiResource(apicode, data, apimsg))

// }

type Response struct {
	Code int
	Msg  string
	Data string
}

func run() {
	e := echo.New()
	// e.POST("/v1/StsToken", StsToken)
	// e.POST("/v2/StsToken", StsToken_v2)
	// e.POST("/v2/StsToken_v1", ApiTest)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{200, "success.echo", "ok"})
	})
	e.Logger.Fatal(e.Start(":80"))
}

// func Redis_() string {
// 	conn, err := redis.Dial("tcp", "192.168.248.126:6379")
// 	if err != nil {
// 		fmt.Println("connect redis error :", err)
// 		return "connect redis error"
// 	}
// 	defer conn.Close()
// 	name, err := redis.String(conn.Do("GET", "name"))
// 	if err != nil {
// 	}
// 	if len(name) > 0 {

// 	} else {
// 		times := strconv.FormatInt(time.Now().Unix(), 10)
// 		conn.Do("SET", "name", times)
// 		conn.Do("expire", "name", 50)
// 		fmt.Println("set redis")
// 		name = times
// 	}
// 	return name
// }

func main() {
	run()
}
