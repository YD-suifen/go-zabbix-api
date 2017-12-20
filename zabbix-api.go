package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

//一下所有的数据类型定义，可以去看zabbixApi文档，

type AuthAPI struct {          //这里的思想是看了人家写好的一些库赖操作的，然后这里首先定义好用户的信息认证，
	url string
	yonghu string
	mima string
	id int
	auth string
}

type Request struct {            //这里定义了发起请求所所需的数据格式，这里的params参数也是类型定义为所有类型，这样就是对之后的 各种请求
                                   //都可以很好的统一使用。
	Jsonrpc string `json:"jsonrpc"`
	Method string `json:"methoed"`
	Params interface{} `json:"params"`
	Auth string `json:"auth"`
	Id int `json:"id"`
}

type Response struct {                   //这里也是一样定义了请求结果返回的数据类型，这里的result也是定位了所有类型，道理一样。
	Jsonrpc string `json:"jsonrpc"`
	Result interface{} `json:"result"`
	Id int `json:"id"`
}

//这里就是跟我们定义的用户类写请求方法。
func (req *AuthAPI) Requestinfo(methon2 string, data2 interface{}) (data Response, web string) {
	//因为咱们这最后可定是持续监听的，所以这里的id做了一下，下一次请求加1.
	id := req.id
	req.id = req.id + 1
	//初始化请求类型
	reqjson := Request{"2.0", methon2, data2, req.auth, id}
	reqjsonzhuhuan, err := json.Marshal(reqjson)   //解析为json数据类型
	if err != nil {
		fmt.Println(err)
	}
	//建立http请求，方式为post
	requestjsoninfo, err2 := http.NewRequest("POST", req.url, bytes.NewReader(reqjsonzhuhuan))
	if err2 != nil {
		fmt.Println(err2)
	}
	//设置都被为application/json，官方的
	requestjsoninfo.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err3 := client.Do(requestjsoninfo)
	if err3 != nil {
		fmt.Println(err3)
	}
	body ,err4 := ioutil.ReadAll(resp.Body)   //读取请求返回的数据

	var xindefanhui Response //定义一个返回数据类型

	jsonerr := json.Unmarshal([]byte(body), &xindefanhui)  //解析返回数据json，到定义的数据类型
	if jsonerr != nil {
		fmt.Println(jsonerr)
	}

	if err4 != nil {
		fmt.Println(err4)
	}
	return xindefanhui, string(body)  // 我这里多打印了一个字符串类型的。作为页面展示
}

//定义登录数据类型。我这里是吧登录的请求跟其他的分开了，因为试了好几次不会弄，auth参数的问题，所以，登陆的数据类型个整个请求单独进行
type loginjiegouti struct {
	Jsonrpc string `json:"jsonrpc"`
	Method string `json:"method"`
	Params loginparams `json:"params"`
	Id int `json:"id"`
}
type loginparams struct {
	User string `json:"user"`
	Password string `json:"password"`
}
func (auth2 *AuthAPI) Login() {    //这是登陆的方法
	loginjsonzhuan := loginjiegouti{Jsonrpc:"2.0", Method:"user.login", Params:loginparams{auth2.yonghu,auth2.mima}, Id:1}
	loginjson, err := json.Marshal(loginjsonzhuan)
	if err != nil {
		fmt.Println(err)
	}

	loginjsoninfo, err2 := http.NewRequest("POST", auth2.url, bytes.NewReader(loginjson))
	if err2 != nil {
		fmt.Println(err2)
	}
	loginjsoninfo.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err3 := client.Do(loginjsoninfo)
	if err3 != nil {
		fmt.Println(err3)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var loginfanhui Response

	jsonerr := json.Unmarshal(body, &loginfanhui)
	if jsonerr != nil {
		fmt.Println(jsonerr)
	}
    auth2.auth = loginfanhui.Result.(string)  //这里就是得到了auth
}


//这是报警1的数据类型
type baojinjsonjiegou struct {
	Output []string `json:"output"`
	Filter baojinfilter `json:"filter"`
	Sortfied string `json:"sortfied"`
	Sortorder string `json:"sortfier"`

}
type baojinfilter struct {
	Value int `json:"value"`
}


func (info *AuthAPI) Baojininfo() string {
	var s baojinjsonjiegou
	s.Output = append(s.Output, "triggerid")
	s.Output = append(s.Output, "description")
	s.Output = append(s.Output, "priority")
	s.Filter.Value = 1
	s.Sortfied = "priority"
	s.Sortorder = "DESC"
	baojindata, baojininfofanhui := info.Requestinfo("trigger.get", s)
	fmt.Println(baojindata)
	return baojininfofanhui
}

type baojin2jiegou struct {
	Triggerids string `json:"triggerids"`
	Output string `json:"output"`
	SelectFunctions string `json:"select_functions"`

}

func (req *AuthAPI) baojin2() string  {
	var bb baojin2jiegou
	bb.Triggerids = "16322"
	bb.Output = "extend"
	bb.SelectFunctions = "extend"
	_, baojin2webinfo := req.Requestinfo("trigger.get", bb)
	//fmt.Println(baojin2data)
	return baojin2webinfo

}

//这是获取主机列表的数据类型
type Hostpeizhijiegou struct {
	Output []string `json:"output"`
	SelectInterfaces []string `json:"select_interfaces"`
}
//方法
func (req *AuthAPI) Hostconf() (hostdata Response, hostwebinfo string)  {
	var params Hostpeizhijiegou
	params.Output = append(params.Output, "hostid")
	params.Output = append(params.Output, "host")
	params.SelectInterfaces = append(params.SelectInterfaces, "interfaceid")
	params.SelectInterfaces = append(params.SelectInterfaces, "ip")
	hostconfdata, hostconfwebinfo := req.Requestinfo("host.get", params)
	fmt.Println(hostconfdata.Result)
	fmt.Println("llllllllll",hostconfdata)
	fmt.Println("sssssssssss",hostconfdata.Jsonrpc)
	fmt.Println("sssssssssss",hostconfdata.Id)
	//for _, v := range hostconfdata.Result{
	//	fmt.Println(string(v))
	//}
	//m := hostconfdata.Result.(map[string]interface{})
	////fmt.Println(m)
	//for k, v := range m{
	//	switch vv := v.(type) {
	//	case string:
	//		fmt.Println(k,vv)
	//	case int:
	//		fmt.Println(k, "is int", vv)
	//	case float64:
	//		fmt.Println(k,"is float64",vv)
	//	case []interface{}:
	//		fmt.Println(k, "is an array:")
	//		for i, u := range vv {
	//			fmt.Println(i,u)
	//
	//		}
	//
	//	}
	//}

	return hostconfdata, hostconfwebinfo

}




func main()  {
	//初始化
	nihao := AuthAPI{url:"http://jk.szytest.com/api_jsonrpc.php",yonghu:"fanchuncheng", mima:"Fanchuncheng11", id:1}
	nihao.Login()
	//启动一个http服务打印，这是自己添加的，测试使用，
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		web1 := nihao.Baojininfo()
		fmt.Fprintln(writer, web1)
	})
	http.HandleFunc("/baojin2", func(writer http.ResponseWriter, request *http.Request) {
		web2 := nihao.baojin2()
		fmt.Fprintln(writer, web2)
	})
	http.HandleFunc("/hostlist", func(writer http.ResponseWriter, request *http.Request) {
		webdata, _ := nihao.Hostconf()
		//func dayin(){
		//	m := webdata.Result.(map[string]interface{})
		//	//fmt.Println(m)
		//	for k, v := range m{
		//		switch vv := v.(type) {
		//		case string:
		//			fmt.Println(k,vv)
		//		case int:
		//			fmt.Println(k, "is int", vv)
		//		case float64:
		//			fmt.Println(k,"is float64",vv)
		//		case []interface{}:
		//			fmt.Println(k, "is an array:")
		//			for i, u := range vv {
		//				fmt.Println(i,u)
		//
		//			}
		//
		//		}
		//	}
		//
		//}



		fmt.Fprintln(writer, webdata)
	})
	http.ListenAndServe("127.0.0.1:8000", nil)
}