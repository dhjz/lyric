package dao

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

var DefaultTimeFormat = "2006-01-02 15:04:05"

type BaseBean struct {
	Page  int `json:"page" gorm:"-:all"`
	Limit int `json:"limit" gorm:"-:all"`
}

func init() {
	Db = CreateDb()
}

func CreateDb() *gorm.DB {
	// exePath, err := os.Executable()
	// if err != nil {
	// 	panic("failed to get executable path")
	// }
	// dbPath := filepath.Join(filepath.Dir(exePath), "dlrc.db")
	// fmt.Println("Connecting to database at:", dbPath)

	_db, err := gorm.Open(sqlite.Open("./dlrc.db"), &gorm.Config{
		Logger: func() logger.Interface {
			// return logger.Default.LogMode(logger.Info)
			return logger.Default.LogMode(logger.Silent)
		}(),
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	fmt.Println("connect database dlrc.db...")

	// 迁移 schema
	_db.AutoMigrate(&Lyric{})

	return _db
}

// MarshalAndGzip 将任意类型的数据 marshal 为 JSON，然后 gzip 压缩
func Compress(data interface{}) []byte {
	// 1. 将任意类型的数据序列化为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed to marshal data to JSON: %w \n", err)
		return nil
	}

	// 2. 使用 gzip 压缩 JSON 数据
	var compressedBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBuffer)
	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		gzipWriter.Close() // 确保关闭 writer
		fmt.Printf("failed to gzip compress data: %w\n", err)
		return nil
	}
	err = gzipWriter.Close() // 必须关闭 writer，才能 flush
	if err != nil {
		fmt.Printf("failed to close gzip writer: %w\n", err)
		return nil
	}

	return compressedBuffer.Bytes()
}

// UngzipAndUnmarshal 将 gzip 压缩的 []byte 数据解压缩，然后 unmarshal 为指定类型的对象
// target 必须是一个指向你想要反序列化到的类型的指针，例如 &LyricResult{}
func Decompress(compressedData []byte, target interface{}) error {
	if len(compressedData) == 0 {
		return fmt.Errorf("input compressedData is empty")
	}

	// 1. 使用 gzip 解压缩数据
	compressedBuffer := bytes.NewReader(compressedData)
	gzipReader, err := gzip.NewReader(compressedBuffer)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return fmt.Errorf("failed to decompress data: %w", err)
	}

	// 2. 将解压缩后的 JSON 数据反序列化为 target 指向的对象
	err = json.Unmarshal(decompressedData, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal decompressed data: %w", err)
	}

	return nil
}
