package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(2000, 10000) // 每秒最多处理 10 个请求，令牌桶大小为 1
const (
	DBUser     = "root"
	DBPassword = "mbssmbss"
	DBName     = "sample"
)

var index = 0
var kvServerURLs = []string{"http://kv_server0:8081", "http://kv_server1:8081", "http://kv_server2:8081"}

var db *sql.DB

func main() {
	time.Sleep(25 * time.Second)
	r := gin.Default()
	initDB()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // 允许所有源访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Use(limitMiddleware)

	r.GET("/put", proxyHandler)
	r.GET("/get", proxyHandler)
	r.GET("/cmd", proxyHandler)
	r.GET("/db", dbHandler)
	fmt.Println("Starting api server on :8080")
	r.Run(":8080")
}

func limitMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func proxyHandler(c *gin.Context) {
	var resp *http.Response
	var err error

	client := &http.Client{Timeout: 10 * time.Second}

	kvserverurl := kvServerURLs[index]
	index = (index + 1) % len(kvServerURLs)

	switch c.Request.URL.Path {
	case "/put":
		resp, err = client.Get(kvserverurl + "/put?" + c.Request.URL.RawQuery)
	case "/get":
		resp, err = client.Get(kvserverurl + "/get?" + c.Request.URL.RawQuery)
	case "/cmd":
		resp, err = client.Get(kvserverurl + "/cmd?" + c.Request.URL.RawQuery)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func initDB() {

	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DBUser, DBPassword, DBName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// 验证连接
	err = db.Ping()
	if err != nil {
		panic("failed to connect database")
	}
}

func dbHandler(c *gin.Context) {
	const query = "SELECT id, `key`, `value` FROM rqyc_kv LIMIT 100"

	rows, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	fmt.Println("rows: ", rows)
	keyValues := make([]map[string]interface{}, 0)

	// Scan all rows into memory to minimize network communication
	cols, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		kv := make(map[string]interface{})
		for i, col := range cols {
			kv[col] = string(rawResult[i])
		}

		keyValues = append(keyValues, kv)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keyValues)
}
