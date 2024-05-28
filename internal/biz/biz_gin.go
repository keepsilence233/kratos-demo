package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tealeg/xlsx"
	"io"
	"kratos-demo/internal/common"
	"net/http"
	"strconv"
	"time"
)

const (
	UserTableName = "account"

	// 批量创建时的batchSize,一次批量create的数量
	BatchCreateUserSize = 2
)

// DB的结构 Notice 实际业务中如果用redis查询，可以用数据库兜底～
type Account struct {
	Id               uint32 `gorm:"primaryKey;column:id" json:"id" structs:"id"`
	Name             string `gorm:"column:name" json:"name" structs:"name" validate:"required,min=1"`
	Gender           string `gorm:"column:desc" json:"desc" structs:"desc"`
	IdCard           string
	Province         string
	City             string
	Region           string
	Company          string
	SocialCreditCode string
	Phone            string
	CreateAt         time.Time `gorm:"column:create_at;default:null" json:"createAt" structs:"create_at"`
	UpdateAt         time.Time `gorm:"column:update_at;default:null" json:"updateAt" structs:"update_at"`
}

func (u Account) TableName() string {
	return UserTableName
}

type BizGinRepo interface {
	BatchCreateExcelDataToDBAndRedis(ctx context.Context, users []*Account) error

	ListFruits(ctx context.Context) (any, error)
}

// GinUsecase的方法
func (gc *GinUsecase) SaveExcelData(ctx context.Context, xlsxFile *excelize.File) error {

	// 测试数据
	//userLst := gc.genTestUserLst()

	//解析excel的数据
	userLst, _, lxRrr := gc.readExcel(xlsxFile)
	if lxRrr != nil {
		return errors.New(fmt.Sprintf("解析excel文件失败err: %v", lxRrr))
	}

	// 写入数据库与redis
	errBatchCreate := gc.ginRepo.BatchCreateExcelDataToDBAndRedis(ctx, userLst)
	if errBatchCreate != nil {
		return errBatchCreate
	}

	return nil
}

func (gc *GinUsecase) SaveExcelAccountData(c *gin.Context, xlsxFile *excelize.File) error {

	saveData, errData, err := gc.readExcel(xlsxFile)

	if err != nil {
		return errors.New(fmt.Sprintf("解析excel文件失败err: %v", err))
	}

	if errData != nil {
		errDataStr, _ := json.Marshal(saveData)
		fmt.Printf(fmt.Sprintf("解析成功异常数据: %v", string(errDataStr)))
		err := ToExcel([]string{`序号`, `姓名`, `性别`, `身份证号码`,
			`省`, `市`, `县`, `单位`, `社会信用代码`, `电话`}, errData,
			"/Users/chinese.youth/github/kratos/kratos-demo/errorData"+time.Now().String()+".xlsx")
		if err != nil {
			return err
		}
	}

	saveDataStr, err := json.Marshal(saveData)
	fmt.Printf(fmt.Sprintf("解析成功保存账户数据: %v", string(saveDataStr)))

	//resp, err := http.Get("http://127.0.0.1:8080/api/student/getStudent/1")
	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)

	if len(saveData) <= 100 {
		err = gc.callXBaseCreateAccount(saveData)
		if err != nil {
			return err
		}
	} else {
		splitAccounts := common.SplitSlice(saveData, 100)
		for _, saveData := range splitAccounts {
			err = gc.callXBaseCreateAccount(saveData)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// readExcel 读取excel
func (gc *GinUsecase) readExcel(xlsx *excelize.File) ([]*Account, []*Account, error) {
	// 根据名字获取cells的内容，返回的是一个[][]string
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

	var saveAccounts []*Account
	var errAccounts []*Account
	for i, row := range rows {
		// Notice 去掉第一行是excel表头部分
		if i == 0 {
			continue
		}

		data := gc.parseRow(row)

		if data.Name == "" || data.IdCard == "" || data.Phone == "" {
			errAccounts = append(errAccounts, &data)
			continue
		}

		saveAccounts = append(saveAccounts, &data)
	}
	return saveAccounts, errAccounts, nil
}

func (gc *GinUsecase) parseRow(row []string) Account {
	return Account{
		Id:               cast.ToUint32(row[0]),
		Name:             row[1],
		Gender:           row[2],
		IdCard:           row[3],
		Province:         row[4],
		City:             row[5],
		Region:           row[6],
		Company:          row[7],
		SocialCreditCode: row[8],
		Phone:            row[9],
		CreateAt:         time.Now(),
		UpdateAt:         time.Now(),
	}
}
func (gc *GinUsecase) callXBaseCreateAccount(accounts []*Account) error {

	_, err := json.Marshal(accounts)
	if err != nil {
		return err
	}
	var data = "{\n    \"Name\":\"wangwu1\",\n    \"Info\":\"wangwu1-info2\",\n    \"Status\":\"123424\"\n}"
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/student/createStudent", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf(fmt.Sprintf("请求创建用户接口返回信息: %v", string(body)))

	return nil
}

// 生成io.ReadSeeker  参数 titleList 为Excel表头，dataList 为数据
func ToExcel(titleList []string, dataList []*Account, filePath string) error {
	// 生成一个新的文件
	file := xlsx.NewFile()
	// 添加sheet页
	sheet, _ := file.AddSheet("Sheet1")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
	}
	// 插入内容
	for _, data := range dataList {
		row := sheet.AddRow()
		//row.WriteStruct(v, -1)
		row.AddCell().Value = strconv.Itoa(int(data.Id))
		row.AddCell().Value = data.Name
		row.AddCell().Value = data.Gender
		row.AddCell().Value = data.IdCard
		row.AddCell().Value = data.Province
		row.AddCell().Value = data.City
		row.AddCell().Value = data.Region
		row.AddCell().Value = data.Company
		row.AddCell().Value = data.SocialCreditCode
		row.AddCell().Value = data.Phone
	}

	//var buffer bytes.Buffer
	//_ = file.Write(&buffer)
	//content = bytes.NewReader(buffer.Bytes())

	// 将文件写入指定路径
	err := file.Save(filePath)
	if err != nil {
		return err
	}
	return nil
}
