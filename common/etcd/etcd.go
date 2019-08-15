package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
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

// 删除任务
func (jobMgr *JobMgr) DelJob(name string) (oldJob *Job, err error) {
	var (
		jobKey string
		delResp *clientv3.DeleteResponse
		oldJobObj Job
	)

	jobKey = "/cron/jobs/" + name

	// 从etcd中删除它
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 解析一下旧值, 返回它
		if err =json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return
}

// 获取列表
func (jobMgr *JobMgr) ListJobs() (jobList []*Job, err error) {
	var (
		dirKey string
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		job *Job
	)

	// 任务保存的目录
	dirKey = "/cron/jobs/"

	// 获取目录下所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 初始化数组空间
	jobList = make([]*Job, 0)

	// 遍历所有任务, 进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &Job{}
		if err =json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

// 强制杀死进程
func (jobMgr *JobMgr) KillJob(name string) (err error) {
	// 更新一下key=/cron/killer/任务名
	var (
		killerKey string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
	)

	// 通知worker杀死对应任务
	killerKey = "/cron/jobs/" + name

	// 让worker监听到一次put操作, 创建一个租约让其稍后自动过期即可
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	// 租约ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记
	if _, err = jobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
