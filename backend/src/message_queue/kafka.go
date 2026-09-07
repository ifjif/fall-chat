package messagequeue

import (
	"context"
	"strconv"
	"sync"
	"time"
	"zzc/fall-chat/src/utils"

	"github.com/segmentio/kafka-go"
)

const (
	// 重试次数的 header key
	HeaderRetryCount = "x-retry-count"
	// 最大重试次数
	MaxRetryCount = 3
)

type MessageHandleFunc func(ctx context.Context, msg *kafka.Message, commitFunc func()) error

type ConsumerGroup struct {
	brokers []string
	topic   string
	groupID string
	// 可有多个hanlder来消费同一个 消息
	handlers  []MessageHandleFunc
	readers   []*kafka.Reader
	dlqWriter *kafka.Writer // 死信队列
	mu        sync.Mutex
}

// 创建一个消费者组
func NewConsumerGroup(brokers []string, topic, groupID string) *ConsumerGroup {
	return &ConsumerGroup{
		brokers:  brokers,
		topic:    topic,
		groupID:  groupID,
		handlers: make([]MessageHandleFunc, 0),
		readers:  make([]*kafka.Reader, 0),
		// 初始化 dlq writer
		dlqWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic + "-dlq",
		}),
	}
}

// 为消费者组注册一个处理器
func (cg *ConsumerGroup) RegisterHandler(handler MessageHandleFunc) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.handlers = append(cg.handlers, handler)
}

// 启动一个该消费者组的 raeder实例
// 调用多次可以启动多个实例（消费者组模式）
func (cg *ConsumerGroup) StartInstance() {
	cg.mu.Lock()
	defer cg.mu.Unlock()

	reader := getkafkaReader(cg.brokers, cg.topic, cg.groupID)
	cg.readers = append(cg.readers, reader)

	go cg.consumeLoop(reader)
}

// 单个reader的消费循环
func (cg *ConsumerGroup) consumeLoop(reader *kafka.Reader) {
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			utils.LogWarn("获取消息失败 topic[%s] group[%s]: %v", cg.topic, cg.groupID, err)
			break
		}

		go cg.processMessage(reader, &msg)
	}
}

// 处理单条消息，会调用组中所有handler
func (cg *ConsumerGroup) processMessage(reader *kafka.Reader, msg *kafka.Message) {
	cg.mu.Lock()
	handlers := make([]MessageHandleFunc, len(cg.handlers))
	copy(handlers, cg.handlers)
	cg.mu.Unlock()

	if len(handlers) == 0 {
		utils.LogWarn("消费者组 topic[%s] group[%s] 未注册任何处理器，不消费", msg.Topic, cg.groupID)
		return
	}

	// 设置提交函数
	var once sync.Once
	commitFunc := func() {
		once.Do(func() {
			// 提交
			if err := reader.CommitMessages(context.Background(), *msg); err != nil {
				utils.LogWarn("提交offset 失败: topic[%s] err[%v]", msg.Topic, err)
			}
		})
	}

	// 退出前尝试提交，无论业务层是否提交
	defer commitFunc()

	// 调用所有处理器
	for _, hanlder := range handlers {
		err := hanlder(context.Background(), msg, commitFunc)
		if err != nil {
			utils.LogError("处理器处理消息失败 topic[%s] group[%s]: %v", msg.Topic, cg.groupID, err)

			// 1.获取当前重试次数
			retryCount := cg.getRetryCount(msg)

			// 2.判断是否到达最大重试次数
			if retryCount >= MaxRetryCount {
				utils.LogWarn("消息到达最大重试次数 [%d], 转入死信队列(DLQ)", MaxRetryCount)
				cg.sendToDLQ(msg, err)
				// 转入dlq后，提交offset, 避免毒丸消息阻塞队列
				reader.CommitMessages(context.Background(), *msg)
				return
			}

			// 3.未达上限，增加重试次数，重新写回 Topic
			cg.sendToRetry(msg, retryCount+1)
			// 原消息进行提交
			reader.CommitMessages(context.Background(), *msg)
			return
		}
	}
}

// 从header 中解析重试次数
func (cg *ConsumerGroup) getRetryCount(msg *kafka.Message) int {
	for _, h := range msg.Headers {
		if h.Key == HeaderRetryCount {
			count, _ := strconv.Atoi(string(h.Value))
			return count
		}
	}

	return 0
}

// 将消息重新发回原topic, 并携带递增的重试次数
func (cg *ConsumerGroup) sendToRetry(msg *kafka.Message, retryCount int) {
	// 复制消息，更新header
	newHeaders := make([]kafka.Header, 0, len(msg.Headers))
	for _, h := range msg.Headers {
		if h.Key != HeaderRetryCount {
			newHeaders = append(newHeaders, h)
		}
	}

	newHeaders = append(newHeaders, kafka.Header{
		Key:   HeaderRetryCount,
		Value: []byte(strconv.Itoa(retryCount)),
	})

	err := kw.WriteMessages(context.Background(), kafka.Message{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: newHeaders,
	})
	if err != nil {
		utils.LogError("消息重试发送失败: topic[%s] err[%v]", msg.Topic, err)
	} else {
		utils.LogInfo("消息已重新发送进行重试：topic[%s], retryCount[%d]", msg.Topic, retryCount)
	}
}

// 将消息发送到死信队列，并附加错误信息
func (cg *ConsumerGroup) sendToDLQ(msg *kafka.Message, originalErr error) {
	// 在header中记录失败原因
	dlqHeaders := make([]kafka.Header, len(msg.Headers))
	copy(dlqHeaders, msg.Headers)
	dlqHeaders = append(dlqHeaders, kafka.Header{
		Key:   "x-dlq-reason",
		Value: []byte(originalErr.Error()),
	})

	err := cg.dlqWriter.WriteMessages(context.Background(), kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: dlqHeaders,
	})
	if err != nil {
		utils.LogError("写入死信队列(DLQ)失败: topic[%s] err[%v]", msg.Topic, err)
	} else {
		utils.LogInfo("消息写入死信队列(DLQ): topic[%s]", msg.Topic)
	}
}

type Dispatcher struct {
	mu     sync.RWMutex
	groups map[string]*ConsumerGroup // key: topic:groupID
}

var (
	kw         *kafka.Writer
	kr         *kafka.Reader
	dispatcher = &Dispatcher{
		groups: make(map[string]*ConsumerGroup),
	}
)

func GetDispatcher() *Dispatcher {
	return dispatcher
}

func KWriter() *kafka.Writer {
	return kw
}

func kafkaInit() {
	kafkaInitGlobalWriter()
}

func kafkaInitGlobalWriter() {
	writerConfig := getWriterConfig()
	kw = kafka.NewWriter(*writerConfig)
	utils.LogInfo("kafka [global writer]: address[%s]", writerConfig.Brokers[0])
}

func getWriterConfig() *kafka.WriterConfig {
	return &kafka.WriterConfig{
		Brokers: []string{"kafka39:9092"},
	}
}

func (d *Dispatcher) GetOrCreateGroup(topic, groupID string) *ConsumerGroup {
	brokers := []string{"kafka39:9092"}
	key := topic + ":" + groupID
	d.mu.RLock()
	group, exists := d.groups[key]
	d.mu.RUnlock()

	if exists {
		return group
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	// 双检查
	if group, exists := d.groups[key]; exists {
		return group
	}

	group = NewConsumerGroup(brokers, topic, groupID)
	d.groups[key] = group
	return group
}

// 便捷方法，注册处理器并启动一个实例
func (d *Dispatcher) RegisterHandlerAndStartInstance(topic, groupID string, handler MessageHandleFunc) {
	group := d.GetOrCreateGroup(topic, groupID)
	group.RegisterHandler(handler)
	group.StartInstance()
}

// 启动一个专门消费死信队列的reader
// dlqTopic: 死信队列 Topic名称（如：chat-topic-dlq）
// handler: 接受到死信队列后的处理逻辑(通常用于告警或记录日志)
func (d *Dispatcher) StartDLQConsumer(dlqTopic, groupID string, handler MessageHandleFunc) {
	brokers := []string{"kafka39:9092"}
	reader := getkafkaReader(brokers, dlqTopic, groupID)

	go d.dlqConsumeLoop(reader, dlqTopic, handler)
}

func (d *Dispatcher) dlqConsumeLoop(reader *kafka.Reader, topic string, handler MessageHandleFunc) {
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			utils.LogWarn("DLQ 获取消息失败：topic[%s] err[%v]", topic, err)
			break
		}

		go d.processDLQMessage(reader, msg, handler)
	}
}

func (d *Dispatcher) processDLQMessage(reader *kafka.Reader, msg kafka.Message, handler MessageHandleFunc) {
	var once sync.Once
	commitFunc := func() {
		once.Do(func() {
			// 提交
			if err := reader.CommitMessages(context.Background(), msg); err != nil {
				utils.LogWarn("DLQ 提交offset失败: topic[%s] err[%v]", msg.Topic, err)
			}
		})
	}

	defer commitFunc()

	if err := handler(context.Background(), &msg, commitFunc); err != nil {
		utils.LogError("DLQ Handler 处理失败：topic[%s] err[%v]", msg.Topic, err)
	}
}

// 将死信队列中的消息重新发回 主topic
func RepublishToMainTopic(mainTopic string, msg *kafka.Message) error {
	// 清除之前重试和DLQ留下的header
	var cleanHeaders []kafka.Header
	for _, h := range msg.Headers {
		if h.Key != HeaderRetryCount && h.Key != "x-dlq-reason" {
			cleanHeaders = append(cleanHeaders, h)
		}
	}

	return kw.WriteMessages(context.Background(), kafka.Message{
		Topic:   mainTopic,
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: cleanHeaders,
	})
}

func getkafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	readerConfig := getReaderConfig(brokers, topic, groupID)
	kr = kafka.NewReader(*readerConfig)
	utils.LogInfo("kafka [reader]: address[%s] topic[%s] group[%s]",
		readerConfig.Brokers[0], readerConfig.Topic, readerConfig.GroupID)
	return kr
}

func getReaderConfig(brokers []string, topic, groupID string) *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Brokers:           brokers,
		Topic:             topic,
		GroupID:           groupID,
		MinBytes:          10e3,
		MaxBytes:          10e6,
		SessionTimeout:    30 * time.Second,
		HeartbeatInterval: 10 * time.Second,
		RebalanceTimeout:  30 * time.Second,
	}
}

/*
*=============================== 使用示例1 =======================================
*
	go func() {
		kw := messagequeue.KWriter()
		timeout := time.NewTicker(2 * time.Second)
		age := 12
		for {
			<-timeout.C
			utils.LogDebug("开始发送kafka消息")
			data, _ := json.Marshal(map[string]any{
				"name":    "张三",
				"age":     age,
				"type":    "group",
				"content": "这是content",
			})
			kw.WriteMessages(context.Background(), kafka.Message{
				Topic: "test-topic-1",
				Key:   []byte("test-topic-key-1"),
				Value: data,
			})
			kw.WriteMessages(context.Background(), kafka.Message{
				Topic: "test-topic-2",
				Key:   []byte("test-topic-key-2"),
				Value: data,
			})
			kw.WriteMessages(context.Background(), kafka.Message{
				Topic: "test-topic-3",
				Key:   []byte("test-topic-key-3"),
				Value: data,
			})
			age += 1
		}
	}()

	testGroup1 := messagequeue.GetDispatcher().GetOrCreateGroup("test-topic-1", "test-topic-group-1")
	testGroup1.RegisterHandler(func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogInfo("这是测试组1的hanlder1: %v", v)
		return nil
	})
	testGroup1.RegisterHandler(func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogWarn("这是测试组1的handler2: %v", v)
		return nil
	})
	testGroup1.StartInstance()

	testGroup2 := messagequeue.GetDispatcher().GetOrCreateGroup("test-topic-2", "test-topic-group-2")
	testGroup2.RegisterHandler(func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogError("这是测试组2的handler1: %v", v)
		return nil
	})
	testGroup2.StartInstance()
	testGroup2.StartInstance()
	testGroup2.StartInstance()

	messagequeue.GetDispatcher().RegisterHandlerAndStartInstance("test-topic-3", "test-group-3", func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogDebug("这是测试组3的hanlder1: %v", v)
		return nil
	})
*
*
*
*
*=============================== 使用示例2 =======================================
*
	go func() {
		kw := messagequeue.KWriter()
		timeout := time.NewTicker(2 * time.Second)
		age := 12
		for {
			<-timeout.C
			utils.LogDebug("开始发送kafka消息")
			data, _ := json.Marshal(map[string]any{
				"name":    "张三",
				"age":     age,
				"type":    "group",
				"content": "这是content",
			})
			kw.WriteMessages(context.Background(), kafka.Message{
				Topic: "test-dlq-1",
				Key:   []byte("test-dlq-key-1"),
				Value: data,
			})
			age += 1
		}
	}()

	testDLQGroup1 := messagequeue.GetDispatcher().GetOrCreateGroup("test-dlq-1", "test-dlq-group-1")
	testDLQGroup1.RegisterHandler(func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogError("接受到testDLQGroup1的数据：%v 进行错误返回", v)
		return fmt.Errorf("这是测试错误")
	})
	testDLQGroup1.StartInstance()

	messagequeue.GetDispatcher().StartDLQConsumer("test-dlq-1-dlq", "test-dlq-1-dlq-group", func(ctx context.Context, msg *kafka.Message,commitFunc func()) error {
		var v map[string]any
		json.Unmarshal(msg.Value, &v)
		utils.LogDebug("收到test-dlq-1-dlq的消息: %v", v)

		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
		v["re"] = "这是重新推送的"
		data, _ := json.Marshal(v)
		messagequeue.RepublishToMainTopic(msg.Topic, &kafka.Message{
			Headers: msg.Headers,
			Key:     msg.Key,
			Value:   data,
		})

		return nil
	})
*
*
*
*
*
*
*
*
*
*/
