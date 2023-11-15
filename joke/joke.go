package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func calculateHash(input string) string {
	// 创建一个 SHA256 哈希对象
	hasher := sha256.New()

	// 将字符串转换为字节数组并写入哈希对象
	hasher.Write([]byte(input))

	// 计算哈希值并转换为十六进制字符串
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString
}

// joke model
type Joke struct {
	ID         uint   `gorm:"primaryKey"`
	Context    string `gorm:"type:text"`
	UpdateTime string `gorm:"type:datetime"`
	HashValue  string `gorm:"type:varchar(32);uniqueIndex"`
}

// 定义一个结构体，用于映射JSON数据
type Result struct {
	StatusCode string    `json:"statusCode"`
	Desc       string    `json:"desc"`
	ResultList []Message `json:"result"`
}

// 定义一个结构体，用于映射JSON中的消息数据
type Message struct {
	ID         int    `json:"id"`
	Content    string `json:"content"`
	UpdateTime string `json:"updateTime"`
}

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("mydatabase.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("无法连接到数据库")
	}

	// 自动迁移模型，确保表存在
	db.AutoMigrate(&Joke{})

	body, err := request()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 定义一个Result结构体实例
	var result Result

	// 使用json.Unmarshal解析JSON数据到Result结构体
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	// 打印解析后的结果
	fmt.Printf("StatusCode: %s\nDesc: %s\n", result.StatusCode, result.Desc)

	// 遍历解析后的消息列表
	for _, msg := range result.ResultList {
		fmt.Printf("ID: %d, Content: %s, UpdateTime: %s\n", msg.ID, msg.Content, msg.UpdateTime)
		// 计算哈希值
		hashValue := calculateHash(msg.Content)
		// 创建Joke实例
		joke := Joke{Context: msg.Content, UpdateTime: msg.UpdateTime, HashValue: hashValue}
		// 将Joke实例写入数据库
		db.Create(&joke)
	}

}

func request() ([]byte, error) {
	uri := "https://eolink.o.apispace.com/xhdq/common/joke/getJokesByRandom"

	payload := url.Values{}
	payload.Set("pageSize", "5")
	req, _ := http.NewRequest("POST", uri, strings.NewReader(payload.Encode()))

	req.Header.Add("X-APISpace-Token", os.Getenv("X_APISPACE_TOKEN"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
