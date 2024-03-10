package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

// UserFile: 用户文件结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
		"`file_size`,`upload_at`)values(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	//返回数据库变化的记录数量
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

// QueryUserFileMetas: 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name = ? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
