package service_gin

import "github.com/gin-gonic/gin"

func RegisterHTTPServer(us *GinService) *gin.Engine {
	router := gin.New()

	rootGrp := router.Group("/api")

	// 用户相关API
	accountGrp := rootGrp.Group("/account")
	// helloWorld
	accountGrp.GET("/sayhi", us.helloKratosGin)
	// 上传excel文件存入DB与Redis中
	accountGrp.POST("/uploadAccount", us.UploadExcelAccounts)
	// 使用gin下载
	accountGrp.GET("/downloadFileGin", us.DownloadFileGin)

	return router
}
