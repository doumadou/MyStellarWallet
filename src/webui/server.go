package webui

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	//"encoding/json"
	"encoding/json"
	"stellar"
	"time"
)

type Asset struct {
	Balance string
	Limit   string
	Issuer  string
	Code    string
}

type User struct {
	UserName       string
	Password       string
	NativeBalances string
	AssetInfo      []Asset
}

type SendPostData struct {
	Sendto string
	Amount string
}

type OfferInfo struct {
	Code   string
	Issuer string
	Price  string
	Amount string
	ID     int
	Type   string
}

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

type walletController struct {
}

var sessionMgr *SessionMgr = nil //session管理器

func (this *walletController) CreatedAction(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles("static/created")
	if err != nil {
		fmt.Println(err)
	}

	//username := "testAccount"
	//password := "soojomoon"
	username, password := stellar.CreateAccount()

	t.Execute(w, &User{UserName: username, Password: password})
}

func createdAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/created")
	if err != nil {
		fmt.Println(err)
	}

	//username := "testAccount"
	//password := "soojomoon"
	username, password := stellar.CreateAccount()

	t.Execute(w, &User{UserName: username, Password: password})
}

func logoutAction(w http.ResponseWriter, r *http.Request) {

	sessionMgr.EndSession(w, r)

	http.Redirect(w, r, "/index", http.StatusFound)
}

func indexAction(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t, err := template.ParseFiles("static/index")
		if err != nil {
			fmt.Println(err)
		}

		t.Execute(w, nil)
	} else if r.Method == "POST" {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()

		//可以使用template.HTMLEscapeString()来避免用户进行js注入
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := &User{UserName: username, Password: password}

		var sessionID = sessionMgr.StartSession(w, r)
		sessionMgr.SetSessionVal(sessionID, "UserInfo", user)
		http.Redirect(w, r, "/wallet", http.StatusFound)
	}
}

func (this *walletController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {

	if r.Method == "GET" {
		t, err := template.ParseFiles("static/index")
		if err != nil {
			fmt.Println(err)
		}

		t.Execute(w, &User{UserName: user, Password: ""})
	} else if r.Method == "POST" {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()

		//可以使用template.HTMLEscapeString()来避免用户进行js注入
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := &User{UserName: username, Password: password}

		var sessionID = sessionMgr.StartSession(w, r)
		sessionMgr.SetSessionVal(sessionID, "UserInfo", user)
		http.Redirect(w, r, "/wallet", http.StatusFound)
	}
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func sendAction(w http.ResponseWriter, r *http.Request) {

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user == nil || !v {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	if r.Method == "POST" {

		//fmt.Println(r)
		//fmt.Println(r.ParseForm())
		//fmt.Println(r.PostForm)
		r.ParseForm() //解析参数，默认是不会解析的
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		////未知类型的推荐处理方法

		//var f interface{}
		//json.Unmarshal(result, &f)
		//m := f.(map[string]interface{})
		//for k, v := range m {
		//	switch vv := v.(type) {
		//	case string:
		//		fmt.Println(k, "is string", vv)
		//	case int:
		//		fmt.Println(k, "is int", vv)
		//	case float64:
		//		fmt.Println(k, "is float64", vv)
		//	case []interface{}:
		//		fmt.Println(k, "is an array:")
		//		for i, u := range vv {
		//			fmt.Println(i, u)
		//		}
		//	default:
		//		fmt.Println(k, "is of a type I don't know how to handle")
		//	}
		//}

		////结构已知，解析到结构体

		var p SendPostData
		json.Unmarshal(result, &p)
		fmt.Println(p.Sendto)
		fmt.Println(p.Amount)

		err := stellar.SendCurreny(userinfo.Password, p.Sendto, p.Amount)
		if err != nil {
			fmt.Println(err)
			OutputJson(w, -1, "发送失败", userinfo)
		} else {
			OutputJson(w, 0, "success", userinfo)
		}

		return
	}

	t, err := template.ParseFiles("static/send")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(userinfo)

	account, err := stellar.GetAccoutInfo(userinfo.UserName)
	if err != nil {
		fmt.Println(err)
	} else {
		userinfo.NativeBalances = account.GetNativeBalance()
	}

	fmt.Println("Send Pages")
	fmt.Println(userinfo)
	fmt.Println(w)

	t.Execute(w, &userinfo)
}

func tradeAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/trade")
	if err != nil {
		fmt.Println(err)
	}

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	fmt.Println(userinfo)
	t.Execute(w, &userinfo)
}

func (this *walletController) TradeAction(w http.ResponseWriter, r *http.Request, userinfo *User, issuer string, code string) {

	t, err := template.ParseFiles("static/trade")
	if err != nil {
		fmt.Println(err)
	}

	// offers, err := stellar.GetBookOffer(issuer, code)

	t.Execute(w, map[string]interface{}{"user": &userinfo, "issuer": issuer, "code": code})
}

func trustAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/trust")
	if err != nil {
		fmt.Println(err)
	}

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	if r.Method == "POST" {

		r.ParseForm()
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		var p Asset
		json.Unmarshal(result, &p)
		fmt.Println(p.Limit)
		fmt.Println(p.Code)
		fmt.Println(p.Issuer)

		var err error
		if p.Limit == "0" {
			err = stellar.RemoveTrustIssue(userinfo.Password, p.Issuer, p.Code)
			if err != nil {
				fmt.Println(err)
				OutputJson(w, -1, "发送失败", userinfo)
			} else {
				var i int
				for i = 0; i < len(userinfo.AssetInfo); i++ {
					if userinfo.AssetInfo[i].Code == p.Code && userinfo.AssetInfo[i].Issuer == p.Issuer {
						break
					}
				}
				if i > 0 && i < len(userinfo.AssetInfo) {
					userinfo.AssetInfo = append(userinfo.AssetInfo[:i], userinfo.AssetInfo[i+1:]...)
				}
				OutputJson(w, 0, "success", userinfo)
			}
		} else {
			err = stellar.TrustIssue(userinfo.Password, p.Issuer, p.Code)
			if err != nil {
				fmt.Println(err)
				OutputJson(w, -1, "发送失败", userinfo)
			} else {
				userinfo.AssetInfo = append(userinfo.AssetInfo, Asset{"0.0000000", "100000000000.000000", p.Issuer, p.Code})
				OutputJson(w, 0, "success", userinfo)
			}
		}

		return
	}

	fmt.Println(userinfo)
	t.Execute(w, &userinfo)
}

func opsAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/hist_ops")
	if err != nil {
		fmt.Println(err)
	}

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	fmt.Println(userinfo)

	records, err := stellar.GetHistroyOps(userinfo.UserName)
	if err != nil {
		fmt.Println(err)
		t.Execute(w, map[string]interface{}{"user": &userinfo, "error": "error"})
	} else {
		t.Execute(w, map[string]interface{}{"user": &userinfo, "records": records})
	}
}

func paymentsAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/hist_payments")
	if err != nil {
		fmt.Println(err)
	}

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	fmt.Println(userinfo)

	records, err := stellar.GetHistroyPayments(userinfo.UserName)
	if err != nil {
		fmt.Println(err)
		t.Execute(w, map[string]interface{}{"user": &userinfo, "error": "error"})
	} else {
		t.Execute(w, map[string]interface{}{"user": &userinfo, "records": records})
	}
}

func walletAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/wallet")
	if err != nil {
		fmt.Println(err)
	}

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	//userinfo := &User{UserName: "GDYTWJAQDHU6SWHFYRJP7Q6SF6DKFFUI54YFDDR2N7K77G74YIZ3KDQX", Password: "SC4DDNUB2562RZRQ3QGN3JONHFI7XVAE5U4AW5QTQSOJ43NUSFFRZQJ2"}

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)
	///t.Execute(w, &User{username, ""})

	account, err := stellar.GetAccoutInfo(userinfo.UserName)
	if err != nil {
		fmt.Println(err)
	} else {
		userinfo.NativeBalances = account.GetNativeBalance()
		userinfo.AssetInfo = []Asset{}
		for _, balance := range account.Balances {
			if balance.Asset.Type != "native" {
				userinfo.AssetInfo = append(userinfo.AssetInfo, Asset{balance.Balance, balance.Limit, balance.Asset.Issuer, balance.Asset.Code})
			}
		}
	}

	fmt.Println(account)
	fmt.Println(userinfo)
	t.Execute(w, &userinfo)
}

func (this *walletController) WalletAction(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles("static/wallet")
	if err != nil {
		fmt.Println(err)
	}
	//username := "ooooooooo"
	//if username != "" {
	//	test_session_valid(w, r)
	//} else {
	//	//http.Redirect(w, r, "/index", http.StatusFound)
	//}
	t.Execute(w, &User{UserName: user, Password: ""})
}

func (this *walletController) LogoutAction(w http.ResponseWriter, r *http.Request) {
	sessionMgr.EndSession(w, r) //用户退出时删除对应session
	http.Redirect(w, r, "/index", http.StatusFound)
	return
}

func test_session_valid(w http.ResponseWriter, r *http.Request) string {
	var sessionID = sessionMgr.CheckCookieValid(w, r)

	if sessionID == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return sessionID
	}

	return sessionID
}

func walletHandler(w http.ResponseWriter, r *http.Request) {
	// 获取cookie
	// cookie, err := r.Cookie("wallet_name")
	// if err != nil || cookie.Value == ""{
	//     http.Redirect(w, r, "/login/index", http.StatusFound)
	// }

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	userinfo := user.(*User)

	fmt.Println(userinfo)

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	issuer := ""
	code := ""
	if len(parts) > 0 {
		action = strings.Title(parts[0]) + "Action"
		if len(parts) == 3 {
			issuer = parts[1]
			code = parts[2]
		}
	}

	if len(parts) != 3 && len(userinfo.AssetInfo) > 0 {
		issuer = userinfo.AssetInfo[0].Issuer
		code = userinfo.AssetInfo[0].Code
	}

	fmt.Println("pathInfo: ", pathInfo)
	fmt.Println("Action: ", action)

	wallet := &walletController{}
	controller := reflect.ValueOf(wallet)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("Trade") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	issuer_t := reflect.ValueOf(issuer)
	code_t := reflect.ValueOf(code)
	user_t := reflect.ValueOf(userinfo)
	fmt.Println(issuer_t)
	fmt.Println(code_t)
	fmt.Println(user_t)
	method.Call([]reflect.Value{responseValue, requestValue, user_t, issuer_t, code_t})

	fmt.Println("method.Call end!")
}

func refreshBookAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Printf("%s\n", result)

	var p OfferInfo
	json.Unmarshal(result, &p)
	fmt.Println(p.Issuer)
	fmt.Println(p.Code)

	offers, err := stellar.GetBookOffer(p.Issuer, p.Code)
	if err != nil {
		OutputJson(w, -1, "获取委单失败", "")
	} else {
		OutputJson(w, 0, "success", offers)
	}

}

func refreshOfferAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Printf("%s\n", result)

	offers, err := stellar.GetAccoutOffers(string(result))
	if err != nil {
		OutputJson(w, -1, "获取委单失败", "")
	} else {
		OutputJson(w, 0, "success", offers)
	}
}

func createOfferAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Printf("%s\n", result)

	var p OfferInfo
	json.Unmarshal(result, &p)
	fmt.Println(p)

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		OutputJson(w, -1, "末登陆", "")
	}
	userinfo := user.(*User)

	err := stellar.CreateOffer(userinfo.Password, p.Issuer, p.Code, p.Price, p.Amount, p.Type)
	if err != nil {
		OutputJson(w, -1, "下单失败", "")
	} else {
		OutputJson(w, 0, "success", "")
	}
}

func cancelOfferAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Printf("%s\n", result)

	var p OfferInfo
	json.Unmarshal(result, &p)
	fmt.Println(p)

	sessionID := test_session_valid(w, r)

	user, v := sessionMgr.GetSessionVal(sessionID, "UserInfo")

	if user != nil && v {
	} else {
		OutputJson(w, -1, "末登陆", "")
	}
	userinfo := user.(*User)

	err := stellar.CancelOffer(userinfo.Password, p.Issuer, p.Code, p.Price, p.ID)
	if err != nil {
		OutputJson(w, -1, "取消失败", "")
	} else {
		OutputJson(w, 0, "success", "")
	}
}

func webui() {

	sessionMgr = NewSessionMgr("TestCookieName", 3600)

	http.Handle("/css/", http.FileServer(http.Dir("static")))
	http.Handle("/js/", http.FileServer(http.Dir("static")))

	//http.HandleFunc("/index", walletHandler)
	//http.HandleFunc("/", walletHandler)
	http.HandleFunc("/index", indexAction)
	http.HandleFunc("/logout", logoutAction)
	http.HandleFunc("/created", createdAction)
	http.HandleFunc("/wallet", walletAction)
	http.HandleFunc("/send", sendAction)
	http.HandleFunc("/trust", trustAction)
	http.HandleFunc("/trade/", walletHandler)
	http.HandleFunc("/trade", walletHandler)
	http.HandleFunc("/hist_payments", paymentsAction)
	http.HandleFunc("/hist_trades", opsAction)
	http.HandleFunc("/refreshBook", refreshBookAction)
	http.HandleFunc("/refreshOffer", refreshOfferAction)
	http.HandleFunc("/createOffer", createOfferAction)
	http.HandleFunc("/cancelOffer", cancelOfferAction)
	http.HandleFunc("/", indexAction)

	//http.ListenAndServe(":8888", Handler: http.TimeoutHandler(http.DefaultServeMux, time.Second*2, "服务端操作超时")

	server := http.Server{
		Addr:    ":8888",
		Handler: http.TimeoutHandler(http.DefaultServeMux, time.Second*60, "服务端操作超时"),
	}

	fmt.Printf("start web server: 0.0.0.0:8888\n")
	server.ListenAndServe()
}

func RunServer() {
	webui()
}
