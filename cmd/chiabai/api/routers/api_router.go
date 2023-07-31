package routers

import (
	// "fmt"
	"fmt"
	"net/http"

	// "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/api/controllers"
	// c_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"


	// "github.com/meta-node-blockchain/meta-node/cmd/chiabai/core"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
)

// SetupRouter sets up the API routes and returns the Gin router.
func InitRouter() *gin.Engine {
	server := controllers.Server{}
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)

	server.Init(cConfig)
	// cli := serverapp.ConnectionHandler()
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	v1 := router.Group("/api/v1")
	//http://localhost:3000/api/v1/test/template/
	{
		v1.StaticFS("", http.Dir("frontend/public"))
		// CreateNFTRoutes(v1)
	}
	router.GET("/ws", func(c *gin.Context) {
		server.WebsocketHandler(c.Writer, c.Request)
	})
	// fmt.Println("server is running on port 2000")

	return router
}

//localhost:2000/api/v1/chiabai/template/
