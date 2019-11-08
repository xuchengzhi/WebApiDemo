package main

import "fmt"
import "github.com/gin-gonic/gin"

import "math/rand"
import "net/http"
import "strconv"
import "time"
import "encoding/base64"
import "encoding/hex"
import "crypto/sha256"
import "crypto/hmac"

// import "github.com/garyburd/redigo/redis"

import "encoding/json"

func HelloPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to bgops,please visit https://xxbandy.github.io!",
	})
}

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
	// fmt.Println(string(key[ml/kl]))
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

func StsToken_v2(c *gin.Context) {
	data := make(map[string]string)
	sts_info := c.PostForm("sts_info")
	key := "xczceshi"

	if len(sts_info) == 0 {
		data["result"] = "5"
		data["msg"] = "error"
		data["info"] = "sts_info not in params"
	} else {

		msg := XorDecodeStr(sts_info, key)
		now := (time.Now().UnixNano() / 1e6)
		var jsons StsJson
		if err := json.Unmarshal([]byte(msg), &jsons); err == nil {
			t, errs := strconv.Atoi(jsons.Times)
			if errs != nil {
			}
			tt := int64(t)
			if now-tt > 1000000 {
				data["result"] = "1"
				data["msg"] = "error"
				data["info"] = "time out"
			} else {

				tmp := HMAC_SHA256(jsons.Times+jsons.Types+jsons.Sys, "qqprivatekey")
				tokens := BASE64EncodeStr(tmp)

				if tokens == jsons.Token {
					toks := "ok" //Redis_()
					data["result"] = "0"
					data["msg"] = "success"
					data["info"] = toks
				} else {
					data["result"] = "2"
					data["msg"] = "error"
					data["info"] = "json 校验失败"
				}
			}
		} else {
			data["result"] = "4"
			data["msg"] = "error"
			data["info"] = "Xor Msg Error"
		}
	}
	var XorMs string
	jsonStr, err := json.Marshal(data)
	if err != nil {
	}
	if is_xor {
		XorMs = XorEncodeStr(string(jsonStr), key)
	} else {
		XorMs = string(jsonStr)
	}
	// fmt.Println(len(XorMs))
	c.String(http.StatusOK, XorMs)
	// fmt.Println(string(jsonStr))
	// ctx.WriteString(string(XorMs))
}

func StsToken(c *gin.Context) {
	apimsg := ""
	var apicode int
	var data string
	sts_info := c.PostForm("sts_info")
	key := "xczceshi"
	msg := XorDecodeStr(sts_info, key)

	var jsons StsJson
	if err := json.Unmarshal([]byte(msg), &jsons); err == nil {
		tmp := HMAC_SHA256(jsons.Times+jsons.Types+jsons.Sys, "qqprivatekey")
		tokens := BASE64EncodeStr(tmp)

		if tokens == jsons.Token {
			toks := "ok"  //Redis_()
			apimsg = toks //"验证通过"
			apicode = 0
			data = "success"
		}

	} else {
		// ctx.StatusCode(iris.StatusBadRequest)
		apimsg = "json 校验失败"
		apicode = 2
		data = "error"
	}
	fmt.Fprintln(gin.DefaultWriter, apicode)
	c.JSON(http.StatusOK, ApiResource(apicode, data, apimsg))

}

func ApiTest(c *gin.Context) {
	times := time.Now()
	c.JSON(http.StatusOK, ApiResource(0, times, "success"))
}

type Response struct {
	Code int
	Msg  string
	Data string
}

func run() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/StsToken", StsToken_v2)
		v1.POST("/StsToken_v1", ApiTest)
		v1.GET("/app/test", func(c *gin.Context) {
			// name := c.Param("name")
			c.String(http.StatusOK, "ok")
		})
		v1.GET("/hello/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "Hello %s", name)
		})

		v1.GET("/line", func(c *gin.Context) {
			// 注意:在前后端分离过程中，需要注意跨域问题，因此需要设置请求头
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			legendData := []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
			xAxisData := []int{120, 240, rand.Intn(500), rand.Intn(500), 150, 230, 180}
			c.JSON(200, gin.H{
				"legend_data": legendData,
				"xAxis_data":  xAxisData,
			})
		})
	}
	//定义默认路由
	r.NoRoute(func(c *gin.Context) {
		// c.JSON(http.StatusNotFound, gin.H{
		// 	"status": 404,
		// 	"error":  "404, page not exists!",
		// })
		c.JSON(http.StatusOK, Response{200, "success gin", "ok"})
	})
	r.Run(":80")
}

// func Redis_() string {
// 	conn, err := redis.Dial("tcp", "192.168.248.126:6379")
// 	if err != nil {
// 		// fmt.Println("connect redis error :", err)
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
// 		conn.Do("expire", "name", 5)
// 		fmt.Println("set redis")
// 		name = times
// 	}
// 	return name
// }

func main() {
	run()
}
