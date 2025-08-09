package tests

import (
	"bytes"
	"log"
	"testing"
	"time"
	"gospring/logging"
	"github.com/stretchr/testify/assert"
)

// 测试用的事件
type TestEvent struct {
	message string
}

func (e *TestEvent) String() string {
	return e.message
}

// 测试用的自定义日志器
type TestLogger struct {
	events []logging.Event
}

func (l *TestLogger) LogEvent(event logging.Event) {
	l.events = append(l.events, event)
}

func (l *TestLogger) GetEvents() []logging.Event {
	return l.events
}

func (l *TestLogger) Clear() {
	l.events = nil
}

// TestNopLogger 测试空日志器
func TestNopLogger(t *testing.T) {
	logger := logging.NopLogger
	event := &TestEvent{message: "test message"}
	
	// NopLogger 应该不做任何事情
	logger.LogEvent(event)
	
	// 测试通过，因为 NopLogger 不会抛出异常
	assert.NotNil(t, logger)
}

// TestConsoleLogger 测试控制台日志器
func TestConsoleLogger(t *testing.T) {
	logger := logging.NewConsoleLogger()
	event := &TestEvent{message: "test console message"}
	
	// ConsoleLogger 应该能够记录事件
	logger.LogEvent(event)
	
	// 测试通过，因为 ConsoleLogger 不会抛出异常
	assert.NotNil(t, logger)
}

// TestStandardLogger 测试标准日志器
func TestStandardLogger(t *testing.T) {
	var buf bytes.Buffer
	stdLogger := log.New(&buf, "", 0)
	logger := logging.NewStandardLogger(stdLogger)
	event := &TestEvent{message: "test standard message"}
	
	logger.LogEvent(event)
	
	// 检查日志是否被写入缓冲区
	assert.Contains(t, buf.String(), "test standard message")
}

// TestMultiLogger 测试多重日志器
func TestMultiLogger(t *testing.T) {
	logger1 := &TestLogger{}
	logger2 := &TestLogger{}
	logger3 := &TestLogger{}
	
	multiLogger := logging.NewMultiLogger(logger1, logger2, logger3)
	event := &TestEvent{message: "test multi message"}
	
	multiLogger.LogEvent(event)
	
	// 检查所有日志器都收到了事件
	assert.Len(t, logger1.GetEvents(), 1)
	assert.Len(t, logger2.GetEvents(), 1)
	assert.Len(t, logger3.GetEvents(), 1)
	assert.Equal(t, "test multi message", logger1.GetEvents()[0].String())
	assert.Equal(t, "test multi message", logger2.GetEvents()[0].String())
	assert.Equal(t, "test multi message", logger3.GetEvents()[0].String())
}

// TestFilteredLogger 测试过滤日志器
func TestFilteredLogger(t *testing.T) {
	testLogger := &TestLogger{}
	
	// 创建一个过滤器，只允许包含 "important" 的消息
	filter := func(event logging.Event) bool {
		return event.String() == "important message"
	}
	
	filteredLogger := logging.NewFilteredLogger(testLogger, filter)
	
	// 测试被过滤的消息
	event1 := &TestEvent{message: "normal message"}
	filteredLogger.LogEvent(event1)
	assert.Len(t, testLogger.GetEvents(), 0)
	
	// 测试通过过滤器的消息
	event2 := &TestEvent{message: "important message"}
	filteredLogger.LogEvent(event2)
	assert.Len(t, testLogger.GetEvents(), 1)
	assert.Equal(t, "important message", testLogger.GetEvents()[0].String())
}

// TestLeveledLogger 测试分级日志器
func TestLeveledLogger(t *testing.T) {
	testLogger := &TestLogger{}
	
	// 创建一个 INFO 级别的日志器
	leveledLogger := logging.NewLeveledLogger(testLogger, logging.LogLevelInfo)
	
	// 测试不同级别的事件
	infoEvent := &logging.ComponentCreated{
		Timestamp:     time.Now(),
		ComponentID:   "test",
		ComponentType: "TestService",
		CreationTime:  time.Millisecond,
	}
	
	leveledLogger.LogEvent(infoEvent)
	assert.Len(t, testLogger.GetEvents(), 1)
}

// TestEventTypes 测试各种事件类型
func TestEventTypes(t *testing.T) {
	// 测试 ContainerCreated 事件
	containerCreated := &logging.ContainerCreated{
		Timestamp: time.Now(),
	}
	assert.Contains(t, containerCreated.String(), "Container created")
	
	// 测试 ComponentRegistered 事件
	componentRegistered := &logging.ComponentRegistered{
		Timestamp:     time.Now(),
		ComponentID:   "testComponent",
		ComponentType: "TestService",
		Scope:         "singleton",
	}
	assert.Contains(t, componentRegistered.String(), "Component registered")
	assert.Contains(t, componentRegistered.String(), "testComponent")
	assert.Contains(t, componentRegistered.String(), "TestService")
	assert.Contains(t, componentRegistered.String(), "singleton")
	
	// 测试 DependencyInjected 事件
	dependencyInjected := &logging.DependencyInjected{
		Timestamp:      time.Now(),
		TargetType:     "TestService",
		DependencyType: "TestRepository",
		FieldName:      "repository",
		ByType:         true,
		ByName:         false,
	}
	assert.Contains(t, dependencyInjected.String(), "Dependency injected")
	assert.Contains(t, dependencyInjected.String(), "TestService")
	assert.Contains(t, dependencyInjected.String(), "TestRepository")
	assert.Contains(t, dependencyInjected.String(), "repository")
	
	// 测试 ComponentCreated 事件
	componentCreated := &logging.ComponentCreated{
		Timestamp:     time.Now(),
		ComponentID:   "testComponent",
		ComponentType: "TestService",
		CreationTime:  time.Millisecond,
	}
	assert.Contains(t, componentCreated.String(), "Component created")
	assert.Contains(t, componentCreated.String(), "testComponent")
	assert.Contains(t, componentCreated.String(), "TestService")
	
	// 测试 LifecycleStarting 事件
	lifecycleStarting := &logging.LifecycleStarting{
		Timestamp:     time.Now(),
		ComponentID:   "testComponent",
		ComponentType: "TestService",
		MethodName:    "Init",
	}
	assert.Contains(t, lifecycleStarting.String(), "Lifecycle starting")
	assert.Contains(t, lifecycleStarting.String(), "testComponent")
	assert.Contains(t, lifecycleStarting.String(), "Init")
	
	// 测试 LifecycleStarted 事件
	lifecycleStarted := &logging.LifecycleStarted{
		Timestamp:     time.Now(),
		ComponentID:   "testComponent",
		ComponentType: "TestService",
		MethodName:    "Init",
		Duration:      time.Millisecond,
		Error:         nil,
	}
	assert.Contains(t, lifecycleStarted.String(), "Lifecycle started")
	assert.Contains(t, lifecycleStarted.String(), "testComponent")
	assert.Contains(t, lifecycleStarted.String(), "Init")
	
	// 测试 ContextStarting 事件
	contextStarting := &logging.ContextStarting{
		Timestamp: time.Now(),
	}
	assert.Contains(t, contextStarting.String(), "Application context starting")
	
	// 测试 ContextStarted 事件
	contextStarted := &logging.ContextStarted{
		Timestamp:      time.Now(),
		Duration:       time.Millisecond,
		ComponentCount: 5,
	}
	assert.Contains(t, contextStarted.String(), "Application context started")
	assert.Contains(t, contextStarted.String(), "components: 5")
}

// TestLoggerIntegration 测试日志器集成
func TestLoggerIntegration(t *testing.T) {
	testLogger := &TestLogger{}
	
	// 创建一个复杂的日志器链
	filter := func(event logging.Event) bool {
		return event.String() != "filtered"
	}
	
	filteredLogger := logging.NewFilteredLogger(testLogger, filter)
	multiLogger := logging.NewMultiLogger(filteredLogger, logging.NopLogger)
	
	// 测试事件传递
	event1 := &TestEvent{message: "normal message"}
	multiLogger.LogEvent(event1)
	assert.Len(t, testLogger.GetEvents(), 1)
	
	// 测试过滤
	event2 := &TestEvent{message: "filtered"}
	multiLogger.LogEvent(event2)
	assert.Len(t, testLogger.GetEvents(), 1) // 应该还是 1，因为被过滤了
}