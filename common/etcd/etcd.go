package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"time"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

// 初始化
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"}, // 集群地址
		DialTimeout: 5 * time.Second, // 连接超时

	}
	if client, err = clientv3.New(config); err != nil {
		return err
	}
	// 得到kv、lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 单例
	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
		lease: lease,
	}

	return nil
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *Job) (oldjob *Job, err error) {
	var (
		jobKey string
		jobVal []byte
		putResp *clientv3.PutResponse
		oldJobObj Job
	)
	jobKey = "/cron/jobs/" + job.Name
	if jobVal, err = json.Marshal(job); err != nil {
		return
	}

	// 保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobVal), clientv3.WithPrevKV()); err != nil {
		return
	}
	// 如果是跟新返回旧值
	if putResp.PrevKv != nil {
		// 对旧值进行反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldjob = &oldJobObj

	}
	return
}
