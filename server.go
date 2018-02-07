package main

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
	// "github.com/iikira/BaiduPCS-Go/pcsconfig"
	"github.com/iikira/BaiduPCS-Go/pcscommand"
	"github.com/iikira/BaiduPCS-Go/pcsconfig"
)

func startApiServer(port uint) error{
	fmt.Printf("api server starting...\n")
  // gin.SetMode(gin.ReleaseMode)
  // define api handlers mostly like app.Commands in main.go
  r := gin.Default()
  v1 := r.Group("/api/v1")
  {
    v1.POST("/test", printBody)
    v1.GET("/w/:a/:b", printParams)
    v1.GET("/login/:bduss", loginHandler)
    v1.GET("/su", printBody)
    v1.GET("/logout", logoutHandler)
    v1.GET("/loglist", loglistHandler)
    v1.GET("/quota", quotaHandler)
    v1.GET("/cd/*path", cdHandler)
    v1.GET("/ls/*path", lsHandler)
    v1.GET("/pwd", pwdHandler)
    v1.GET("/meta/*path", metaHandler)
    v1.GET("/rm", printBody)
    v1.GET("/mkdir", printBody)
    v1.GET("/cp", printBody)
    v1.GET("/mv", printBody)
    v1.GET("/download", printBody)
    v1.GET("/upload", printBody)
    v1.GET("/set", printBody)
    v1.GET("/quit", printBody)
  }
  r.Run(fmt.Sprintf(":%d", port))
  return nil
}

func printBody(c *gin.Context) {
        buf := make([]byte, 1024)
        n, _ := c.Request.Body.Read(buf)
        resp := map[string]string{"body": string(buf[0:n])}
        c.JSON(http.StatusOK, resp)
}

func printParams(c *gin.Context) {
  a := c.Param("a")
  b := c.Param("b")
  c.JSON(200, gin.H{
    "a": a,
    "b": b,
  })
}

func loginHandler(c *gin.Context) {
  bduss := c.Param("bduss")
  _, err := pcsconfig.Config.SetBDUSS(bduss, "", "")
  if err != nil {
    c.JSON(500, gin.H{
      "error": err,
    })
  } else {
    c.JSON(200, gin.H{
      "msg": "登录成功",
    })
  }
}

func logoutHandler(c *gin.Context) {
  // 仅退出当前active账号, 与cmd模式行为不同
  uid := pcsconfig.ActiveBaiduUser.UID
  if len(pcsconfig.Config.BaiduUserList) == 0 || uid == 0 {
    c.JSON(500, gin.H{
      "error": "未设置任何百度帐号, 不能退出",
    })
  }
  if !pcsconfig.Config.CheckUIDExist(uid) {
    c.JSON(500, gin.H{
      "error": "退出用户失败, uid 不存在",
    })
  }
  if !pcsconfig.Config.DeleteBaiduUserByUID(uid) {
    c.JSON(500, gin.H{
      "error": "退出用户失败",
    })
  }
  c.JSON(200, gin.H{
    "msg": "退出用户成功",
  })
}

func loglistHandler(c *gin.Context) {
  c.JSON(200, gin.H{
    "active": pcsconfig.ActiveBaiduUser,
    "all": pcsconfig.Config.GetAllBaiduUserInJSON(),
  })
}

func quotaHandler(c *gin.Context) {
  quota, used, err := pcscommand.GetQuota()
  if err != nil {
    c.JSON(500, gin.H{
      "error": err,
    })
		return
	} else {
    c.JSON(200, gin.H{
      "quota": quota,
      "used": used,
    })
  }
}

func cdHandler(c *gin.Context) {
  path := c.Param("path")
  err := pcscommand.ChangeDirectory(path)
  if err != nil {
    c.JSON(500, gin.H{
      "error": err,
    })
		return
	} else {
    lsHandler(c)
  }
}

func lsHandler(c *gin.Context) {
  path := c.Param("path")
  files, summary, err := pcscommand.Ls(path)
  if err != nil {
    c.JSON(500, gin.H{
      "error": err,
    })
		return
	} else {
    c.JSON(200, gin.H{
      "path": path,
      "summary": summary,
      "children": files,
    })
  }
}

func pwdHandler(c *gin.Context) {
  c.JSON(200, gin.H{
    "path": pcsconfig.ActiveBaiduUser.Workdir,
  })
}

func metaHandler(c *gin.Context) {
  path := c.Param("path")
  meta, err := pcscommand.GetMeta(path)
  if err != nil {
    c.JSON(500, gin.H{
      "error": err,
    })
		return
	} else {
    c.JSON(200, gin.H{
      "meta": meta,
    })
  }
}
