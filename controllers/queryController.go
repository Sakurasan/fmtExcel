package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golib/security"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	// "fmt"
)

type QueryController struct {
	beego.Controller
}

type Trans struct {
	service string
	charset string
	mer_id  string
	//res_format string
	version  string
	signtype string
	sign     string
	order_id string
	mer_date string
}
type GetTrans struct {
	Mer_id               string
	Mer_date             string
	Version              string
	Ret_code             string
	Ret_msg              string
	Trade_state          string
	Trade_no             string
	Transfer_date        string
	Amount               string
	Transfer_settle_date string
	Order_id             string
	Fee                  string
	Purpose              string
}

func (this *QueryController) Get() {
	this.TplName = "query.html"
	var slcK []template.HTML
	// template.HTML("<img src=" + imgbase64 + " /><br>")
	slcK = append(slcK, template.HTML("<option></option>"))
	configMap := ReadJsonConfig()
	for n, _ := range configMap {
		slcK = append(slcK, template.HTML("<option>"+n+"</option>"))
		// fmt.Fprintln(w, fmt.Sprintf("<option>%s</option>", n))
	}
	this.Data["Namelist"] = slcK

}

func (this *QueryController) Post() {
	this.TplName = "query.html"
	var TwoData [][]template.HTML
	var slcK []template.HTML
	// template.HTML("<img src=" + imgbase64 + " /><br>")
	slcK = append(slcK, template.HTML("<option></option>"))
	configMap := ReadJsonConfig()
	for n, _ := range configMap {
		slcK = append(slcK, template.HTML("<option>"+n+"</option>"))
		// fmt.Fprintln(w, fmt.Sprintf("<option>%s</option>", n))
	}
	this.Data["Namelist"] = slcK

	mer_id := this.Input().Get("cate")
	order_id := this.Input().Get("order_id")
	orderlist := strings.Split(order_id, ",")
	fmt.Println("Order List->", orderlist, "\n", orderlist[0], orderlist[1])

	datelist := []string{}
	for _, date := range orderlist {
		datelist = append(datelist, date)
	}
	// var mer_date string
	for _, order := range orderlist {
		mer_date := order[:8]

		// mer_date := order_id[:8]

		data_val := make(url.Values, 0)
		data_val.Add("charset", "UTF-8")
		data_val.Add("mer_date", mer_date)
		data_val.Add("mer_id", mer_id)
		data_val.Add("service", "transfer_query")
		data_val.Add("order_id", order)
		data_val.Add("version", "4.0")
		encodedata := data_val.Encode()
		fmt.Println("加密项:" + encodedata)
		data_val.Add("sign_type", "RSA")

		url1 := "https://pay.soopay.net/spay/pay/payservice.do"
		// configMap := ReadJsonConfig()
		fmt.Println("mer_id", configMap[mer_id])
		fmt.Println("\n--------------------------\n", configMap[mer_id], "\n-------------------------\n")
		sign, err := RsaSignBase64([]byte(encodedata), configMap[mer_id])
		if err != nil {
			fmt.Println("商户签名", err)
		}
		fmt.Println("商户签名后：" + sign)
		data_val.Set("sign", sign)
		data := data_val.Encode()

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url1+"?"+data, nil)
		req.Header.Set("Accept-Charset", "UTF-8")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
		out, _ := httputil.DumpRequestOut(req, true)
		fmt.Println("\n请求内容->", string(out))

		resp, _ := client.Do(req)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		var myhtml io.Reader
		myhtml = bytes.NewReader(body)
		// fmt.Println("\n返回结果->", string(body))

		// fmt.Println(string(body))
		doc, err := goquery.NewDocumentFromReader(myhtml)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(doc.Html())
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~")
		dhead := doc.Find("head")
		content, _ := dhead.Find("meta").Attr("content")
		fmt.Println("\nmycontent->", content)

		mycontent, _ := url.ParseQuery(content)
		fmt.Println("\nParseQuery->", mycontent)
		resptrans := GetTrans{}
		fee, err := strconv.ParseFloat(mycontent.Get("fee"), 64)
		feef := fee / 100
		resptrans.Fee = fmt.Sprintf("%.2f", feef)
		resptrans.Mer_id = mycontent.Get("mer_id")
		resptrans.Mer_date = mycontent.Get("mer_date")
		amt, err := strconv.ParseFloat(mycontent.Get("amount"), 64)
		amtf := amt / 100
		if err != nil {
			fmt.Println("金额转换失败")
		}
		resptrans.Amount = fmt.Sprintf("%.2f", amtf) //strconv.FormatFloat(amtf, 'E', -1, 64)
		resptrans.Order_id = mycontent.Get("order_id")
		resptrans.Purpose = mycontent.Get("purpose")
		resptrans.Ret_code = mycontent.Get("ret_code")
		resptrans.Ret_msg = mycontent.Get("ret_msg")
		resptrans.Transfer_date = mycontent.Get("transfer_date")
		resptrans.Trade_state = StateConv(mycontent.Get("trade_state"))
		resptrans.Trade_no = mycontent.Get("trade_no")
		resptrans.Transfer_settle_date = mycontent.Get("transfer_settle_date")
		resptrans.Purpose = mycontent.Get("purpose")
		resptrans.Version = mycontent.Get("version")
		fmt.Println(resptrans)

		t := reflect.TypeOf(resptrans)
		v := reflect.ValueOf(resptrans)
		var tblist []template.HTML
		for k := 0; k < t.NumField(); k++ {
			// fmt.Printf("%s -- %v \n", t.Filed(k).Name, v.Field(k).Interface())
			// fmt.Sprintf("%s = %s",t.Filed(k).Name, v.Field(k).Interface())
			// fmt.Fprintln(w, fmt.Sprintf("<p>%s = %s</p>", Conv(t.Field(k).Name), v.Field(k).Interface()))

			tblist = append(tblist, template.HTML(fmt.Sprintf("<td>%s</td>", v.Field(k).Interface())))
		}
		TwoData = append(TwoData, tblist)
	}

	// template.HTML("<img src=" + imgbase64 + " /><br>")
	// fmt.Println("\nTblist", tblist)
	// this.Data["Tblist"] = []template.HTML{"<td>8180</td>", "<td>20171114</td>", "<td>4.0</td>", "<td>0000</td>", "<td>查询成功</td>", "<td>成功</td>", "<td>1711140931093844</td>", "<td>20171114</td>", "<td>19206.36</td>", "<td>20171114</td>", "<td></td>", "<td>2.88</td>", "<td>T0商户结算</td>"}
	fmt.Println("二维数组->", TwoData)
	// TwoTwo := [][]string{{"123", "456"}, {"abc", "def"}}
	this.Data["Twodata"] = TwoData

}

func RsaSignBase64(data []byte, key string) (string, error) {
	key_path := "etc/" + key
	fmt.Println("key_path ", key_path)
	mer_privatekey, err := security.GetRsaPrivateKey(key_path)
	if err != nil {
		fmt.Println("No priKey")
		return "", err
	}

	return security.RsaSignSha1Base64(mer_privatekey, data)
}

func ReadJsonConfig() map[string]string {
	data, err := ioutil.ReadFile("./conf/config.json")
	if err != nil {
		return nil
	}
	datajson := []byte(data)
	fmt.Println(string(datajson))
	var jsonmap map[string]string
	err = json.Unmarshal(datajson, &jsonmap)
	if err != nil {
		fmt.Println("ReadJsonConfig .Unmarshal Err")
		return nil
	}
	// fmt.Println(jsonmap)
	return jsonmap
}

func StateConv(state string) string {
	switch state {
	case "0":
		return "创建"
	case "1":
		return "支付中"
	case "3":
		return "失败"
	case "4":
		return "成功"
	case "5":
		return "借款中"
	case "6":
		return "借款已受理"
	case "7":
		return "退款中"
	case "11":
		return "待确认"
	case "12":
		return "已冻结,待财务审核"
	case "13":
		return "待解冻,交易失败"
	case "14":
		return "财务已审核，待财务付款"
	case "15":
		return "财务审核失败，交易失败"
	case "16":
		return "受理成功，交易处理中"
	case "17":
		return "交易失败退单中"
	case "18":
		return "交易失败退单成功"
	}
	return ""
}

func Conv(old string) string {
	switch old {
	case "Mer_id":
		return "商户编号" + old
	case "Mer_date":
		return "订单日期" + old
	case "Version":
		return "版本号" + old
	case "Ret_code":
		return "返回码" + old
	case "Ret_msg":
		return "返回信息" + old
	case "Trade_state":
		return "交易状态" + old
	case "Trade_no":
		return "U付交易号" + old
	case "Transfer_date":
		return "付款日期" + old
	case "Amount":
		return "付款金额(元)" + old
	case "Transfer_settle_date":
		return "付款对账日期" + old
	case "Order_id":
		return "商户唯一订单号" + old
	case "Fee":
		return "手续费" + old
	case "Purpose":
		return "备注" + old
	}
	return ""
}
