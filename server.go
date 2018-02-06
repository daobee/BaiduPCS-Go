package main

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
)

func startApiServer(port uint) error{
	fmt.Printf("api server starting...\n")
  // gin.SetMode(gin.ReleaseMode)
  // define api handlers mostly like app.Commands in main.go
  r := gin.Default()
  v1 := r.Group("/api/v1")
  {
    v1.POST("/channel", printBody)
    v1.GET("/w/:a/:b", printParams)
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
