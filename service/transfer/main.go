package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/mq"
	"filestore-server/store/oss"
)

// ProcessTransfer: 更新文件存储信息
func ProcessTransfer(msg []byte) bool {
	// 1.解析msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 2.根据临时存储文件路径，创建文件句柄
	file, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 3.通过文件句柄将文件内容读出来并且上传到OSS
	err = oss.Bucket().PutObject(pubData.DestLocation, bufio.NewReader(file))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 4.更新文件的存储路径到文件表
	if dblayer.UpdateFileLocation(pubData.Filehash, pubData.DestLocation) {
		return true
	}
	return false
}

func main() {
	log.Println("开始监听转移任务队列...")
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)
}
