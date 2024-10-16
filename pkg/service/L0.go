package service

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"github.com/jordan-wright/email"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//// FindL0 find L0 by current user id
//func (my MysqlService) FindL0(r *request.L0) []models.L0 {
//	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "FindL0"))
//	defer span.End()
//	list := make([]models.L0, 0)
//	q := my.Q.Tx.
//		Model(&models.L0{}).
//		Order("created_at DESC").
//		Where("user_id = ?", r.UserId)
//	if r.Status != nil {
//		//q.Where("status = ?", *r.Status)
//	}
//	//desc := strings.TrimSpace(r.Desc)
//	if desc != "" {
//		q.Where("`desc` LIKE ?", fmt.Sprintf("%%%s%%", desc))
//	}
//	my.Q.FindWithPage(q, &r.Page, &list)
//	return list
//}

//	func (my MysqlService) CreateL0(r *request.CreateL0) (err error) {
//		// create leave to db
//		//var l0 models.L0
//		//copier.Copy(&l0, r)
//		// save fsm uuid
//		//l0.L0Id = fsmUuid
//		//leave.Status = models.LevelStatusWaiting
//		//leave.UserId = r.User.Id
//
// 我需要在这里实现：1.我在F:\L0\L0里有若干个.csv文件，我需要将他们的数据扫描到我的数据库中。.csv的列与我的model已经对应。
//
//		my.Q.Tx.Create(&r)
//		return
//	}
//
// ImportL0FromCSV 扫描指定目录下的所有.csv文件，并将数据导入数据库
func (my MysqlService) CreateL0(directoryPath string) error {
	files, err := filepath.Glob(filepath.Join(directoryPath, "*.csv"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err := my.importL0FromSingleCSV(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (my MysqlService) importL0FromSingleCSV(filePath string) error {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		l0 := models.L0{
			//L0Id:                         // 从record中获取对应的值，例如：strconv.ParseUint(record[0], 10, 64)
			SOURCE_CHAIN:                 record[0], // 假设csv文件中的第一列对应SOURCE_CHAIN
			SOURCE_TRANSACTION_HASH:      record[1], // 假设csv文件中的第二列对应SOURCE_TRANSACTION_HASH
			DESTINATION_CHAIN:            record[2], // 假设csv文件中的第三列对应DESTINATION_CHAIN
			DESTINATION_TRANSACTION_HASH: record[3], // 假设csv文件中的第四列对应DESTINATION_TRANSACTION_HASH
			SENDER_WALLET:                record[4], // 假设csv文件中的第五列对应SENDER_WALLET
			//SOURCE_TIMESTAMP_UTC:         time.Parse(record[5], "2006-01-02 15:04:05"),
			PROJECT:           record[6],
			NATIVE_DROP_USD:   record[7],
			STARGATE_SWAP_USD: record[8],
		}
		my.Q.Tx.Create(&l0)
	}

	return nil
}

func (my MysqlService) UpdateL0ById(id uint, r request.UpdateL0, u models.SysUser) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "UpdateL0ById"))
	defer span.End()
	var L0 models.L0
	my.Q.Tx.
		Where("id = ?", id).
		First(&L0)
	if L0.Id == constant.Zero {
		err = gorm.ErrRecordNotFound
		return
	}
	// check edit permission
	err = my.Q.FsmCheckEditLogDetailPermission(req.FsmCheckEditLogDetailPermission{
		//Category:       req.NullUint(global.FsmCategoryL0),
		//Uuid:           L0.L0Id,
		ApprovalRoleId: u.RoleId,
		ApprovalUserId: u.Id,
		Fields:         []string{"desc", "start_time", "end_time"},
	})
	if err != nil {
		return
	}
	// update
	my.Q.UpdateById(id, r, new(models.L0))
	return
}

func (my MysqlService) DeleteL0ByIds(ids []uint, u models.SysUser) (err error) {
	_, span := tracer.Start(my.Q.Ctx, tracing.Name(tracing.Db, "DeleteL0ByIds"))
	defer span.End()
	list := make([]string, 0)
	my.Q.Tx.
		Model(&models.L0{}).
		Where("id IN (?)", ids).
		Pluck("fsm_uuid", &list)
	if len(list) > 0 {
		my.Q.FsmCancelLogByUuids(req.FsmCancelLog{
			ApprovalRoleId: u.RoleId,
			ApprovalUserId: u.Id,
			Uuids:          list,
		})
	}
	err = my.Q.DeleteByIds(ids, new(models.L0))
	return
}

// 获取linea积分
func (my MysqlService) GetLineaPoints() {
	filePath := "F:/L0/myAddr.csv"
	addresses, err := TargetAddr(filePath)
	if err != nil {
		fmt.Println("处理文件时出错:", err)
		return
	}
	// 示例用法
	err = GetSybilFromGit(addresses, "LineaPonts.xlsx")
	if err != nil {
		fmt.Println("获取Sybil报告时出错:", err)
		return
	}

	fmt.Println("记录完成")
}

func TargetAddr(filePath string) ([]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取 CSV 文件
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	// 提取第一列数据
	var addresses []string
	for _, record := range records {
		if len(record) > 0 {
			addresses = append(addresses, record[0])
		}
	}
	return addresses, nil
}
func GetSybilFromGit(queries interface{}, outputFile string) error {
	var queryStrings []string
	// 处理传入的参数，将其转换为字符串集合
	switch queries := queries.(type) {
	case string:
		queryStrings = []string{queries}
	case []string:
		queryStrings = queries
	default:
		return fmt.Errorf("无效的参数类型")
	}
	for _, query := range queryStrings {
		// 构建URL
		url := fmt.Sprintf("https://kx58j6x5me.execute-api.us-east-1.amazonaws.com/linea/getUserPointsSearch?user=%s", query)
		// 发送HTTP GET请求
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// Parse the JSON response
		var dataList []models.LineaPoints
		err = json.Unmarshal(body, &dataList)
		err = appendToExcel(dataList, outputFile)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Data appended to xxx.xls successfully.")
		}
	}
	return nil
}
func appendToExcel(dataList []models.LineaPoints, filePath string) error {
	var f *excelize.File
	var err error

	// Check if file exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, create a new file
		f = excelize.NewFile()

		// 保存文件
		if err := f.SaveAs(filePath); err != nil {
			fmt.Println("保存文件时出错：", err)
			return err
		}
		fmt.Println("文件创建成功！")
	} else {
		// File exists, open the existing file
		f, err = excelize.OpenFile(filePath)
		if err != nil {
			fmt.Println("打开文件时出错：", err)
			return err
		}
	}

	sheetMap := f.GetSheetMap()
	sheetName := sheetMap[1] // 使用索引 1 获取工作表名称
	if sheetName == "" {
		return err
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}
	rowIndex := len(rows) + 1
	for _, data := range dataList {
		row := []interface{}{
			data.RankXP,
			data.UserAddress,
			data.XP,
			data.ALP,
			data.PLP,
			data.EP,
			data.RP,
			data.VP,
			data.BP,
			data.EAFlag,
		}
		for i, value := range row {
			cell := fmt.Sprintf("%s%d", string('A'+i), rowIndex+1)
			f.SetCellFormula(sheetName, cell, `{"number_format": 49}`)
			f.SetCellValue(sheetName, cell, value)
		}
		rowIndex++
	}
	if err := f.SaveAs(filePath); err != nil {
		return err
	}
	return nil
}

// 监控alex是否开放
func (my MysqlService) AlexMoniter(url string) {
	body, _ := GetResponseForAlex(url)
	fmt.Println(string(body))
	value := 2
	if value < 1 {
		toAddress := "巅峰弃剑<1334642655@qq.com>" // Replace with the recipient's email address
		subject := "Alex监控"
		body := "alex开放交易"
		CreateEmailMessage(toAddress, subject, body)
		fmt.Println("Email sent successfully!")
		fmt.Println("Alex开放")
	} else {
		fmt.Println("适合转到链上")
	}
}

// 监控sol JUP的swap价格
func (my MysqlService) SwapAeUsdcUsdtOnSOL(url string) {
	body, _ := GetResponse(url)
	str, _ := GetDateFromJson(body, "outAmount")
	outstr := string(str)
	out, _ := strconv.ParseFloat(string(str), 64)
	str, _ = GetDateFromJson(body, "inAmount")
	instr := string(str)
	in, _ := strconv.ParseFloat(instr, 64)
	// 转换为整数并检查是否大于1
	value := (in / out) * 0.001
	fmt.Println(value)
	if value < 1 {
		toAddress := "巅峰弃剑<1334642655@qq.com>" // Replace with the recipient's email address
		subject := "Swap监控"
		body := "\\r\\n磨损系数：" + strconv.FormatFloat(value, 'f', -1, 64) + "\\r\\nIn:" + instr + "\\r\\nOut" + outstr
		CreateEmailMessage(toAddress, subject, body)
		fmt.Println("Email sent successfully!")
		fmt.Println("适合进交易所")
	} else {
		fmt.Println("适合转到链上")
	}
}

// 获得{a:b}json中指定值
func GetDateFromJson(jsonData []byte, key string) (string, error) {
	var data map[string]interface{}
	// 解析JSON数据
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "", err
	}
	// 从map中获取指定字段的值
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("key '%s' not found or not a string", key)
	}
	return value, nil
}
func (my MysqlService) GetOminiUsdtOp1(matches []interface{}) {
	url := "https://optimistic.etherscan.io/tokenholdings?a=0x233ddece6a96c49ece6ad9ae820690fe62a28975"
	resp, err := http.Get(url)
	if err != nil {
		//resp = "58</span></a></span></div></td><td>0.1-替代</td>"
	}
	//defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("错误body:" + string(body))
	}
	//fmt.Println(string(body))

	// 使用正则表达式提取目标字符串
	r := regexp.MustCompile(`58</span></a></span></div></td><td>(.*?)</td>`)
	match := r.FindStringSubmatch(string(body))
	if len(match) < 2 {
		fmt.Println("No match found")
	}

	// 提取目标字符串
	targetString := match[1]
	fmt.Println(targetString)

	// 裁剪字符串
	//trimmedString := targetString[:strings.Index(targetString, "</td>")]

	// 转换为整数并检查是否大于1
	value, err := strconv.ParseFloat(targetString, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if value > 1 {
		toAddress := "recipient@example.com" // Replace with the recipient's email address
		subject := "取款马上"
		body := "存入" + strconv.FormatFloat(value, 'f', -1, 64)
		CreateEmailMessage(toAddress, subject, body)
		fmt.Println("Email sent successfully!")
		fmt.Println("提取的值大于1")
	} else {
		fmt.Println("提取的值不大于1")
	}
}

func writeToFile(data string, isGreaterThanOne string) error {
	file, err := os.Create("ominUsdtOp.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("Data: " + data + "\n")
	file.WriteString("Is greater than 1: " + isGreaterThanOne + "\n")

	return nil
}

// 发送邮件
func CreateEmailMessage(fromAddress string, subject string, body string) gmail.Message {
	e := email.NewEmail()
	e.From = fromAddress
	//"1871437892@qq.com"
	e.To = []string{"zyj18883502123@gmail.com"}
	e.Subject = subject
	e.Text = []byte(body)
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "1334642655@qq.com", "cmhtvatsahmvbafg", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	return gmail.Message{}
}

// 批量查询太古
// 批量读出文件中地址
func (my MysqlService) GetTaikoAirDrop() {
	//读出需要查询数据”0x26a360e91e9f158e6e47e6cbe14f02b649509948“ =xxx
	filePath := "F:/L0/TikoAirDropAddr.csv"
	RemoveBlankAndDuplicateRows(filePath)
	addresses, err := TargetAddr(filePath)
	if err != nil {
		fmt.Println("处理文件时出错:", err)
		return
	}
	// 查询
	for _, query := range addresses {
		GetTaikoAirDrop(query)
	}
	fmt.Println("记录完成")
}
func GetTaikoAirDrop(address string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://trailblazer.hekla.taiko.xyz/api/address?address=%s", address), nil)
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
	AppendToCSV("F:\\L0\\TikoAirDrop.csv", string(body))
	return body, nil
}

// 模拟非浏览器访问
func GetResponse(urlStr string) ([]byte, error) {
	resp, err := http.Get("https://quote-api.jup.ag/v6/quote?inputMint=DdFPRnccQqLD4zCHrBqdY95D6hvw6PLWp9DEXj1fLCL9&outputMint=Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB&amount=3000000000000&slippageBps=10&computeAutoSlippage=true&swapMode=ExactIn&onlyDirectRoutes=false&asLegacyTransaction=false&maxAccounts=64&experimentalDexes=Jupiter%2520LO")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func GetResponseForAlex(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 模拟浏览器访问网站，返回
func GetChromeRequest(r string) ([]byte, error) {
	//encodedURL := url.QueryEscape(r)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://quote-api.jup.ag/v6/quote?inputMint=DdFPRnccQqLD4zCHrBqdY95D6hvw6PLWp9DEXj1fLCL9&outputMint=Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB&amount=3000000000000&slippageBps=10&computeAutoSlippage=true&swapMode=ExactIn&onlyDirectRoutes=false&asLegacyTransaction=false&maxAccounts=64&experimentalDexes=Jupiter%20LO"), nil)
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
	return body, nil
}

// 追加记录到指定文档
func AppendToCSV(csvFile string, value string) error {
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{value})
	if err != nil {
		return err
	}

	return nil
}

// 去除.csv空白行重复行
func RemoveBlankAndDuplicateRows(csvFile string) error {

	// 打开CSV文件
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取CSV文件内容
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// 去除空白行
	nonBlankRows := make([][]string, 0)
	for _, row := range rows {
		isEmpty := true
		for _, field := range row {
			if strings.TrimSpace(field) != "" {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			nonBlankRows = append(nonBlankRows, row)
		}
	}

	// 去除重复行
	uniqueRows := make([][]string, 0)
	seenRows := make(map[string]bool)
	for _, row := range nonBlankRows {
		rowString := strings.Join(row, ",")
		if !seenRows[rowString] {
			uniqueRows = append(uniqueRows, row)
			seenRows[rowString] = true
		}
	}

	// 创建新的CSV文件并写入去除空白行和重复行后的内容
	newFile, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(uniqueRows)
	if err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// 区块浏览器
func GetTxSwapVelar(address string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hiro.so/extended/v1/tx/mempool?address=SP1Y5YSTAHZ88XYK1VPDH24GY0HPX5J4JECTMY4A1.univ2-router&limit=30&offset=0&unanchored=true"), nil)
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
	AppendToCSV("F:\\L0\\TikoAirDrop.csv", string(body))
	return body, nil
}

// 前端页面监控
func vueSwapMonitor(address string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens"), nil)
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
	AppendToCSV("F:\\L0\\TikoAirDrop.csv", string(body))
	return body, nil
}
