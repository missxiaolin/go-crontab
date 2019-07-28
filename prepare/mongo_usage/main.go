package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"start_time"` // 开始时间
	EndTime int64 `bson:"end_time"` // 结束时间
}

type LogRecord struct {
	JobName string `bson:"jobName"` // 任务名称
	Command string `bson:"command"` // shell 命令
	Err string `bson:"err"` // 错误
	Content string `bson:"content"` // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间
	result *mongo.InsertOneResult
}

func main ()  {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertOneResult
	)

	client, err =mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"));
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	database = client.Database("test");
	collection = database.Collection("test");

	// 插入测试
	record = &LogRecord{
		JobName: "job10",
		Command: "echo hellow",
		Err: "",
		Content: "hellow",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime: time.Now().Unix() + 10,
		},
	}
	result, err = collection.InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println(err)
	}
	objId := result.InsertedID.(obj)
	fmt.Println()
}
