package db

import (
	"fmt"
	"os"
	"time"
	"wxcloudrun-golang/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbInstance *gorm.DB

// Init 初始化数据库
func Init() error {

	source := "%s:%s@tcp(%s)/%s?readTimeout=1500ms&writeTimeout=1500ms&charset=utf8&loc=Local&&parseTime=true"
	user := os.Getenv("MYSQL_USERNAME")
	pwd := os.Getenv("MYSQL_PASSWORD")
	addr := os.Getenv("MYSQL_ADDRESS")
	dataBase := os.Getenv("MYSQL_DATABASE")
	if dataBase == "" {
		dataBase = "golang_demo"
	}
	source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	fmt.Println("start init mysql with ", source)

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		}})
	if err != nil {
		fmt.Println("DB Open error,err=", err.Error())
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("DB Init error,err=", err.Error())
		return err
	}

	// 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxIdleConns(100)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(200)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = db

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&model.UserModel{},
		&model.PostModel{},
		&model.CommentModel{},
		&model.CategoryModel{},
		&model.UserLikeModel{},
	)
	if err != nil {
		fmt.Println("AutoMigrate error,err=", err.Error())
		return err
	}

	// 初始化默认分类数据
	initDefaultCategories(db)

	fmt.Println("finish init mysql with ", source)
	return nil
}

// Get ...
func Get() *gorm.DB {
	return dbInstance
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return dbInstance
}

// initDefaultCategories 初始化默认分类数据
func initDefaultCategories(db *gorm.DB) {
	var count int64
	db.Model(&model.CategoryModel{}).Count(&count)
	if count > 0 {
		return // 如果已有数据，不重复初始化
	}

	categories := []model.CategoryModel{
		{
			Id:          "cat_001",
			Name:        "全部",
			Code:        "all",
			Icon:        "📋",
			Description: "所有分类的帖子",
			Sort:        0,
		},
		{
			Id:          "cat_002",
			Name:        "技术",
			Code:        "tech",
			Icon:        "💻",
			Description: "技术分享、开发经验、编程技巧",
			Sort:        1,
		},
		{
			Id:          "cat_003",
			Name:        "生活",
			Code:        "life",
			Icon:        "🏠",
			Description: "日常生活、心情分享、生活感悟",
			Sort:        2,
		},
		{
			Id:          "cat_004",
			Name:        "美食",
			Code:        "food",
			Icon:        "🍜",
			Description: "美食制作、餐厅推荐、食谱分享",
			Sort:        3,
		},
		{
			Id:          "cat_005",
			Name:        "旅行",
			Code:        "travel",
			Icon:        "✈️",
			Description: "旅行攻略、景点推荐、游记分享",
			Sort:        4,
		},
		{
			Id:          "cat_006",
			Name:        "读书",
			Code:        "book",
			Icon:        "📚",
			Description: "书籍推荐、读书笔记、读后感",
			Sort:        5,
		},
		{
			Id:          "cat_007",
			Name:        "运动",
			Code:        "sport",
			Icon:        "🏃",
			Description: "运动健身、体育赛事、健康生活",
			Sort:        6,
		},
	}

	for _, category := range categories {
		db.Create(&category)
	}
}
