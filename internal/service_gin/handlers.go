package service_gin

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gogf/gf/util/gconv"
	"io"
	"os"
)

type UploadExcelReq struct {
	Username string `json:"username" uri:"username" form:"username" binding:"required,min=1"`
	Password string `json:"password" uri:"password" form:"password" binding:"required,min=1"`
	// Notice binding标签为required的话,传来的age不能为空值：0 / 空字符串
	Age int `json:"age" uri:"age" form:"age" binding:"min=1"`
}

// 使用gin上传excel
func (g *GinService) UploadExcelAccounts(c *gin.Context) {

	// Notice 请求的header设置一下:
	// Content-Type 设置成 multipart/form-data;boundary=10

	// Notice 线上的话，做一下文件大小的限制，最好在前端限制下～必要的话在后台也读一下文件流的大小做限制
	// 或者限制一下用户数量：比如在后面的逻辑中发现用户数超过了1w，直接返回错误
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(500, Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("上传文件错误err: %v", err),
			Message: nil,
		})
		return
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(500, Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("读取excel文件失败err: %v", err),
			Message: nil,
		})
		return
	}

	// Notice 可以使用gin的context了！
	err = g.srv.Gc.SaveExcelAccountData(c, xlsx)
	if err != nil {
		c.JSON(500, Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("保存excel数据失败err: %v", err),
			Message: nil,
		})
		return
	}
	c.JSON(500, Response{
		Code:    OK,
		Reason:  "上传成功!",
		Message: nil,
	})

}

// 1、使用gin框架下载
func (g *GinService) DownloadFileGin(c *gin.Context) {

	// Notice 到时候换成 服务器 或 项目部署所在的docker容器(先把文件COPY进去) 里面文件的绝对路径
	filePath := "/Users/wanghongwei/Documents/my-kratos-demos/kratosGin/IMG_1984.jpeg"

	file, errFile := os.Open(filePath)
	if errFile != nil {
		c.JSON(500, Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("读取文件出错:%v", errFile),
			Message: nil,
		})
		return
	}
	defer file.Close()

	fileContentDisposition := fmt.Sprintf("attachment;filename=%v", "蓝天.jpeg")
	c.Header("Content-Disposition", fileContentDisposition)

	// Notice 1、本地文件的写法
	//c.File(filePath)

	// Notice 2、如果是从oss中读取的二进制流数据，这样写
	body, _ := io.ReadAll(file) // 模拟二进制流数据
	_, errW := c.Writer.Write(body)
	if errW != nil {
		c.JSON(500, Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("下载文件出错:%v", errFile),
			Message: nil,
		})
		return
	}

	return
}

// 2、注意 这里的http用的是kratos的http
func (g *GinService) DownloadFile(w http.ResponseWriter, c *http.Request) {

	// Notice 到时候换成 服务器 或 项目部署所在的docker容器(先把文件COPY进去) 里面文件的绝对路径
	filePath := "/Users/wanghongwei/Documents/my-kratos-demos/kratosGin/IMG_1984.jpeg"

	file, errFile := os.Open(filePath)
	if errFile != nil {
		// Notice 返回错误响应，实际上返回的是 bytes
		errStr, _ := json.Marshal(Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("读取文件出错:%v", errFile),
			Message: nil,
		})
		w.Write(gconv.Bytes(errStr))

		return
	}
	defer file.Close()

	// Notice 必要的操作！
	fileContentDisposition := fmt.Sprintf("attachment;filename=%v", "蓝天.jpeg")
	w.Header().Set("Content-Disposition", fileContentDisposition)

	// Notice 注意，下载文件、下载从其他地方获取的二进制流数据的处理方法不一样！

	// Notice 1、因为读取的是文件,所以需要seek到文件的开始位置，将file Copy 到 w 中即可～
	//file.Seek(0, 0)
	//io.Copy(w, file)

	// Notice 2、实际上下载文件给前端一个二进制流数据就好了
	// Notice 如果是从oss中获取的文件流，直接将上面的方法注释掉，直接下面这样写就好了:
	// Notice 假设从oss中获取到了这样一个http数据流: 假设body是从oss中获取到的文件二进制流数据
	body, _ := io.ReadAll(file) // 模拟二进制流数据
	_, errWrite := w.Write(body)
	if errWrite != nil {
		errStr, _ := json.Marshal(Response{
			Code:    ErrCodeCommon,
			Reason:  fmt.Sprintf("下载文件出错:%v", errFile),
			Message: nil,
		})
		w.Write(gconv.Bytes(errStr))
	}

	return
}

// 简单的helloword
func (g *GinService) helloKratosGin(c *gin.Context) {

	/*
		// Notice 虽然可以这样做，但是还是不建议，既然集成了gin框架，请求与响应就都用gin来处理吧～
		// 可以使用gin解析http请求中的参数，然后调用GreeterService中的方法去处理～
		in := &v1.HelloRequest{}
		// Notice context可以直接用gin的
		g.srv.SayHello(c, in)
	*/

	// 获取参数（GET请求中querystring的参数、POST请求中如果Content-Type是application/json那么 json/form表单的参数都可以用这个方法获取）
	userName := c.DefaultQuery("username", "")

	c.JSON(200, gin.H{
		"msg": "success",
		"data": gin.H{
			"greetings": fmt.Sprintf("hello %v", userName),
		},
	})

}
