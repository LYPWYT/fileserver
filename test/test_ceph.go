package main

import (
	"fmt"
	"gopkg.in/amz.v1/s3"

	"filestore-server/store/ceph"
)

func main() {
	bucket := ceph.GetCephBucket("testbucket1")

	// 创建一个新的bucket对象
	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		fmt.Printf("create bucket err: %v\n", err)
	}

	// 查询这个bucket下面指定条件的object keys
	res, err := bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("object keys: %v\n", res)
	}

	// 新上传一个对象
	err = bucket.Put("/testupload/a.txt", []byte("just for test"), "octet-stream", s3.PublicRead)
	fmt.Printf("upload err: %v", err)

	// 查询这个bucket下面制定条件的object keys
	res, err = bucket.List("", "", "", 100)
	if err != nil {
		fmt.Printf("object keys: %v\n", res)
	}
}
