package views

import (
	"net/http"
	"strconv"

	"github.com/chengjoey/web-terminal/connections"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ShellWs(c *gin.Context) {
	var err error
	msg := c.DefaultQuery("msg", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := connections.Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	sshClient, err := connections.DecodedMsgToSSHClient(msg)
	if err != nil {
		c.Error(err)
		return
	}
	if sshClient.IpAddress == "" || sshClient.Password == "" {
		c.Error(&ApiError{Message: "IP地址或密码不能为空", Code: 400})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
}
