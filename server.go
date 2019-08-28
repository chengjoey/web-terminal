package main

import (
	"net/http"

	"github.com/chengjoey/web-terminal/views"
	"github.com/gin-gonic/gin"
)

//对产生的任何error进行处理
func JSONAppErrorReporter() gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *views.ApiError
			switch err.(type) {
			//如果产生的error是自定义的结构体,转换error,返回自定义的code和msg
			case *views.ApiError:
				parsedError = err.(*views.ApiError)
			default:
				parsedError = &views.ApiError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
			c.IndentedJSON(parsedError.Code, parsedError)
			return
		}

	}
}

//设置所有跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("./templates/*.html")
	server.Static("/static", "./templates/static")
	server.Use(JSONAppErrorReporter())
	server.Use(CORSMiddleware())
	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	server.GET("/ws", views.ShellWs)
	server.Run("0.0.0.0:5001")
}
