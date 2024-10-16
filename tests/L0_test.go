package tests

import (
	"gin-web/models"
	"gin-web/pkg/service"
	"os"
	"testing"
	"time"
)

// 监控omnilending
func TestMysqlService_CreateL0(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	q := service.New(ctx)
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// 每隔5分钟调用一次q.GetOminiUsdtOp1(nil)函数
				q.GetOminiUsdtOp1(nil)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(time.Hour * 9999)
	close(quit)
}

// 临时测试指定方法。
func TestAlexMoniter(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	q := service.New(ctx)
	q.AlexMoniter("https://app.alexlab.co/")
}

// 监控网站：https://quote-api.jup.ag/v6/quote?inputMint=DdFPRnccQqLD4zCHrBqdY95D6hvw6PLWp9DEXj1fLCL9&outputMint=Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB&amount=3000000000000&slippageBps=10&computeAutoSlippage=true&swapMode=ExactIn&onlyDirectRoutes=false&asLegacyTransaction=false&maxAccounts=64&experimentalDexes=Jupiter%20LO
func TestSwapAeUsdcUsdtOnSOL(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	url := "https://quote-api.jup.ag/v6/quote?inputMint=DdFPRnccQqLD4zCHrBqdY95D6hvw6PLWp9DEXj1fLCL9&outputMint=Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB&amount=3000000000000&slippageBps=10&computeAutoSlippage=true&swapMode=ExactIn&onlyDirectRoutes=false&asLegacyTransaction=false&maxAccounts=64&experimentalDexes=Jupiter%20LO\n"
	alex := "https://app.alexlab.co/"
	q := service.New(ctx)
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				q.SwapAeUsdcUsdtOnSOL(url)
				q.AlexMoniter(alex)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(time.Hour * 9999)
	close(quit)
}

// 批量查询taiko
func TestTaikoAirdrop(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	q := service.New(ctx)
	q.GetTaikoAirDrop()
}

// L0 查询
func TestL0(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	//q := service.New(ctx)
	//q.GetL0Sibel2Git()
}

// 存入defi数据到数据库
func TestVelar(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")

	Config()
	Mysql()
	//c := new(gin.Context)
	//v1.CreateVelarTokens(c)
	var re models.Message
	q := service.New(ctx)
	q.CreateVelarToken(&re)
	//获取交易所价格
	//q.GetPriceFromOkx("")
}
func TestYoken(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	//c := new(gin.Context)
	//调取v1操作层，操作层调用serve层
	//v1.CreateLeave(c)
	//v1.CreateVelarTokens(c)
	//q := service.New(ctx)
	//q.GetVelarToken("")
}

// 获取linea积分数量
func TestLinea(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	q := service.New(ctx)
	q.GetLineaPoints()
}

//思路：
//1.获取交易所API，
//2.监控swap价格，
//3.比对交易所价格Difi价格
//4.锁定defi价格
//5.保证defi交易成功，交易所再进行交易。
