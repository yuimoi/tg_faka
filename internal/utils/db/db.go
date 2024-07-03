package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"tg_go_faka/internal/models"
	"tg_go_faka/internal/utils/config"
)

var DB *gorm.DB

func InitDB() {
	var db *gorm.DB
	var err error

	dsn := filepath.Join(config.GetRootDir(), ".env", "db.db")
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	// sqlite库默认行为是不开启外键校验的，要手动开启
	DB.Exec("PRAGMA foreign_keys = ON")

	autoMigrate()
	//generateCode()
}

//func generateCode() {
//	genOutputPath := path.Join(config.GetRootDir(), "internal", "dao")
//	// 初始化Gen对象
//	g := gen.NewGenerator(gen.Config{
//		OutPath: genOutputPath, // 生成的代码输出路径
//		//FieldNullable:  true,          // 生成带有空指针类型的字段
//		//FieldCoverable: true,          // 生成带有零值的字段
//		//FieldSignable:  true,          // 生成带有可签名的字段
//	})
//
//	// 使用GORM数据库对象
//	g.UseDB(DB)
//
//	// 生成所有模型的代码
//	for _, model := range models.MyModels {
//		g.ApplyBasic(model)
//	}
//
//	// 执行代码生成
//	g.Execute()
//}

func autoMigrate() {
	for _, model := range models.MyModels {
		var _ models.MyModel = model //顺便校验接口
		if err := DB.AutoMigrate(model); err != nil {
			panic(err)
		}
	}
}
