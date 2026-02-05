package main

import (
	"fmt"
	"log"
	"os"
	"pr-business/config"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败:", err)
	}
	defer sqlDB.Close()

	fmt.Println("✅ 数据库连接成功")

	// 读取迁移脚本
	script, err := os.ReadFile("migrations/004_rename_current_role_to_active_role.sql")
	if err != nil {
		log.Fatal("读取迁移脚本失败:", err)
	}

	// 执行迁移
	result := db.Exec(string(script))
	if result.Error != nil {
		log.Fatal("执行迁移失败:", result.Error)
	}

	fmt.Println("✅ 迁移执行成功: current_role → active_role")

	// 验证迁移
	var columnName string
	err = db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'active_role'").Scan(&columnName).Error
	if err != nil {
		log.Fatal("验证迁移失败:", err)
	}

	fmt.Printf("✅ 验证成功: 列 %s 已存在\n", columnName)
}
