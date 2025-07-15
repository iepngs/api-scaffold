package main

import (
	"flag"
	"fmt"
	"log"

	"bmstock/internal/router"
	"bmstock/pkg/di"

	"github.com/gofiber/fiber/v2"
)

var (
	// 在编译时通过 -ldflags 注入
	version = "dev" // 默认值
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	showVersion := flag.Bool("v", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("bmstock version %s\n", version)
		return
	}

	// 初始化依赖注入容器
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// 创建 Fiber 应用
	app := fiber.New()

	// 初始化路由
	router.SetupRoutes(app, container)

	// 启动服务
	port := container.Config.App.Port
	log.Printf("Starting bmstock version: %s, on port: %s", version, port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
