package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//func FindVelar(c *gin.Context) {
//	ctx := tracing.RealCtx(c)
//	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "FindVelar"))
//	defer span.End()
//	var r request.Velar
//	req.ShouldBind(c, &r)
//	user := GetCurrentUser(c)
//	r.UserId = user.Id
//	my := service.New(c)
//	list := my.FindVelar(&r)
//	resp.SuccessWithPageData(list, &[]response.Velar{}, r.Page)
//}

// 创建
//
//	func CreateVelarTokens(c *gin.Context) {
//		//ctx := tracing.RealCtx(c)
//		//_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "CreateVelar"))
//		//defer span.End()
//		//client := &http.Client{}
//		json, err := http.Get("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens")
//		c.BindJSON(json)
//		var r request.CreateVelar
//		req.ShouldBind(c, &r)
//		my := service.New(c)
//		err = my.CreateVelarToken(&r)
//		resp.CheckErr(err)
//		resp.Success()
//	}
func CreateVelarTokens(c *gin.Context) {
	//client := &http.Client{}
	//req, err := http.NewRequest("GET", fmt.Sprintf("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens"), nil)
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	//resp, err := client.Do(req)
	//defer resp.Body.Close()

	//var re request.CreateVelar
	////body, err := GetVelarToken("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens")
	////
	////c.BindJSON(body)
	////fmt.Println(body)
	//// 将解析后的 JSON 数据绑定到请求上下文中
	////fmt.Println(string(body))
	////utils.Json2Struct(string(body), &r)
	//c.JSON(http.StatusOK, string(body))
	////c.Bind(string(body))
	//
	//req.ShouldBind(c, &r)
	//my := service.New(c)
	//err = my.CreateVelarToken(&r)
	//resp.CheckErr(err)
	//resp.Success()
	//body, err := GetVelarToken("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens")
	//// 将API的返回值作为JSON响应返回给客户端
	//temp := string(body)
	//fmt.Println(temp)
	////c.Data(http.StatusOK, "application/json", body)
	//
	//c.JSON(http.StatusOK, &temp)
	//req.ShouldBind(c, &re)
	//my := service.New(c)
	//err = my.CreateVelarToken(&re)
	//resp.CheckErr(err)
	//resp.Success()
}

func GetVelarToken(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	return body, err
}

func GetVelarToken1(address string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(address), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	return body, err
}

//// UpdateVelarById
//func UpdateVelarById(c *gin.Context) {
//	ctx := tracing.RealCtx(c)
//	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "UpdateVelarById"))
//	defer span.End()
//	var r request.UpdateVelar
//	req.ShouldBind(c, &r)
//	id := req.UintId(c)
//	my := service.New(c)
//	user := GetCurrentUser(c)
//	err := my.UpdateVelarById(id, r, user)
//	resp.CheckErr(err)
//	resp.Success()
//}
//
//// BatchDeleteVelarByIds
//func BatchDeleteVelarByIds(c *gin.Context) {
//	ctx := tracing.RealCtx(c)
//	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "BatchDeleteVelarByIds"))
//	defer span.End()
//	var r req.Ids
//	req.ShouldBind(c, &r)
//	my := service.New(c)
//	user := GetCurrentUser(c)
//	err := my.DeleteVelarByIds(r.Uints(), user)
//	resp.CheckErr(err)
//	resp.Success()
//}
