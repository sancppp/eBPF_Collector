package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库连接信息
const (
	DBUser     = "root"
	DBPassword = "mbssmbss"
	DBName     = "sample"

	Stress = false
)

// 数据库模型
type KeyValue struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"type:varchar(255);uniqueIndex"`
	Value string `gorm:"type:text"`
}

// TableName 设置 KeyValue 的表名为 "rqyc"
func (KeyValue) TableName() string {
	return "rqyc_kv"
}

var db *gorm.DB

func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DBUser, DBPassword, DBName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get DB from GORM")
	}

	// 配置数据库连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库
	db.AutoMigrate(&KeyValue{})
}

// 处理/put请求
func putHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(w, "Missing key or value parameter", http.StatusBadRequest)
		return
	}
	addCORSHeaders(w)

	var kv KeyValue
	// 尝试查找现有记录
	if err := db.Where("`key` = ?", key).First(&kv).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果记录不存在，插入新记录
			kv = KeyValue{Key: key, Value: value}
			if err := db.Create(&kv).Error; err != nil {
				http.Error(w, "Failed to insert data: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data inserted successfully"))
			return
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 如果记录存在，更新值
	if err := db.Model(&kv).Update("value", value).Error; err != nil {
		http.Error(w, "Failed to update data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data updated successfully"))
}

// 处理/get请求
func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}
	addCORSHeaders(w)

	var kv KeyValue
	if err := db.First(&kv, "`key` = ?", key).Error; err != nil {
		// 返回200状态码，但是数据不存在的提示
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data not found"))
		} else {
			http.Error(w, "Failed to query data: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"key": kv.Key, "value": kv.Value})
}

// 处理/cmd请求
func cmdHandler(w http.ResponseWriter, r *http.Request) {
	cmdStr := r.URL.Query().Get("cmd")
	if cmdStr == "" {
		http.Error(w, "Missing cmd parameter", http.StatusBadRequest)
		return
	}
	addCORSHeaders(w)

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Command execution failed: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func main() {
	time.Sleep(20 * time.Second)
	initDB()

	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/cmd", cmdHandler)

	fmt.Println("Starting kv server on :8081")
	http.ListenAndServe(":8081", nil)
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
