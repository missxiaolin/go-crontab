package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	config clientv3.Config
	client *clientv3.Client
	err error
	kv clientv3.KV
)

func init()  {
	// 客户端配置
	config = clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 用户读写etcd的键值对
	kv = clientv3.NewKV(client)
}

func main() {
	//etcdPut()
	//etcdGet()
	//etcdDel()
}

func etcdDel() {
	var (
		delResp *clientv3.DeleteResponse
	)

	if delResp, err = kv.Delete(context.TODO(), "name"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("getResp", delResp.PrevKvs)
}

func etcdPut() {
	var (
		puResp *clientv3.PutResponse
	)

	if puResp, err = kv.Put(context.TODO(),"name", "xiaolin"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Revision", puResp.Header.Revision)
}

func etcdGet()  {
	var (
		getResp *clientv3.GetResponse
	)

	if getResp, err = kv.Get(context.TODO(), "name"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("getResp", getResp.Kvs)
}

