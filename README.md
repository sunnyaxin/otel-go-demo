# OpenTelemetry Go 演示项目

这是一个使用 GoFrame 框架和 OpenTelemetry 的 Go 应用演示项目，展示了如何集成可观测性功能，包括日志、指标和链路追踪。

## 📌 分支说明

本项目包含两个分支，分别演示不同的 OpenTelemetry 集成方式：

### 🤖 main 分支 - 自动插桩（Auto-Instrumentation）

**特点**：
- ✅ **零代码侵入** - 无需修改应用代码即可获得追踪能力
- ✅ **快速部署** - 使用 eBPF 技术在运行时自动注入追踪
- ✅ **统一管理** - 通过环境变量集中配置可观测性

**技术实现**：
- 使用 `opentelemetry-go-instrumentation` 自动插桩代理
- 通过 Docker 容器的 `pid: "service:app"` 模式注入追踪
- 需要特权模式和 eBPF 支持

**适用场景**：
- 快速为现有应用添加可观测性
- 不想修改业务代码
- 需要统一的可观测性方案
- 概念验证（PoC）阶段

**Docker Compose 服务**：
- `app` - 应用容器
- `otel-collector` - OpenTelemetry Collector
- `otel-go-agent` - 自动插桩代理（关键）

---

### ✍️ manual 分支 - 手动插桩（Manual Instrumentation）

**特点**：
- ✅ **精细控制** - 可以自定义 span 名称、属性和采样策略
- ✅ **业务语义** - 添加业务相关的追踪信息和事件
- ✅ **更好的性能** - 避免运行时注入开销
- ✅ **生产就绪** - 更适合生产环境部署

**技术实现**：
- 直接集成 OpenTelemetry Go SDK
- 使用 `otlphttp` 包初始化追踪导出器
- 代码中显式创建和管理 span
- 不需要额外的自动插桩代理容器

**代码示例**（manual 分支）：
```go
import (
    "github.com/gogf/gf/contrib/trace/otlphttp/v2"
    "go.opentelemetry.io/otel"
)

// 初始化追踪
var shutdown, _ = otlphttp.Init(serviceName, endpoint, path)
defer shutdown(ctx)

// 手动创建 span
func addManualTrace(ctx context.Context) {
    tracer := otel.GetTracerProvider().Tracer("otel-go-demo")
    ctx, span := tracer.Start(ctx, "hello-outer-span")
    span.AddEvent("hello-event")
    traceInner(ctx)
    defer span.End()
}
```

**适用场景**：
- 需要详细的业务追踪信息
- 对性能有极致要求
- 生产环境部署
- 需要深度定制可观测性

**Docker Compose 服务**：
- `app` - 应用容器
- `otel-collector` - OpenTelemetry Collector
- ~~无需自动插桩代理~~

---

### 🔄 如何切换分支

```bash
# 查看所有分支
git branch -a

# 切换到自动插桩分支
git checkout main

# 切换到手动插桩分支
git checkout manual

# 查看当前分支
git branch
```

### 📊 分支对比表

| 特性 | main 分支（自动） | manual 分支（手动） |
|------|------------------|-------------------|
| 代码侵入性 | 无 | 需要修改代码 |
| 部署复杂度 | 较高（需要特权容器） | 较低（标准容器） |
| 追踪粒度 | 粗粒度（HTTP请求级） | 细粒度（可自定义） |
| 性能开销 | 略高（eBPF） | 较低 |
| 学习曲线 | 平缓 | 陡峭 |
| 生产就绪度 | 实验性 | 成熟稳定 |
| 自定义能力 | 有限 | 完全控制 |
| 适用阶段 | 开发/测试 | 测试/生产 |

---

> **💡 建议**：
> - **初学者或快速体验**：从 `main` 分支开始，无需修改代码即可看到效果
> - **生产环境**：使用 `manual` 分支，获得更好的控制和性能
> - **学习路径**：先体验 `main` 分支理解概念 → 再学习 `manual` 分支深入掌握

---

## 项目概述

本项目演示了以下内容：

- 使用 GoFrame v2 构建 Web 应用
- 集成 OpenTelemetry 进行可观测性
- Prometheus 指标导出
- 结构化 JSON 日志（兼容 Google Cloud Logging 格式）
- OpenTelemetry Go 自动插桩
- 使用 Docker Compose 进行容器化部署

## 技术栈

- **Go**: 1.23.0
- **GoFrame**: v2.9.4 - Go 企业级应用开发框架
- **OpenTelemetry**: 可观测性标准
- **Prometheus**: 指标收集和监控
- **Docker**: 容器化部署

## 项目结构

```
.
├── cmd/
│   └── server/
│       ├── main.go                      # 主应用入口
│       └── JsonOutputsForLogger.go      # 自定义 JSON 日志处理器
├── auto/                                # 自动化脚本
│   └── dev                              # 开发环境脚本
├── docker-compose.yml                   # Docker Compose 配置
├── Dockerfile                           # 应用容器镜像
├── otel-collector-config.yaml           # OpenTelemetry Collector 配置
├── go.mod                               # Go 模块依赖
└── go.sum                               # Go 依赖校验
```

## 功能特性

### 1. Web 服务端点

- `GET /hello` - 简单的 Hello World 端点，带有日志和追踪
- `GET /metrics` - Prometheus 指标导出端点

### 2. 可观测性

#### 日志
- 结构化 JSON 日志输出
- 兼容 Google Cloud Logging 格式
- 自动包含 TraceID 和 SpanID 用于关联追踪

#### 指标
- Prometheus 格式的指标导出
- 通过 `/metrics` 端点暴露
- OpenTelemetry Collector 自动采集

#### 链路追踪
- OpenTelemetry Go 自动插桩
- HTTP 请求自动追踪
- 通过 OTLP 协议导出到 Collector

## 快速开始

### 前置要求

- Docker 和 Docker Compose
- Go 1.23.0 或更高版本（用于本地开发）

### 使用 Docker Compose 运行

1. 克隆项目并进入目录：
```bash
cd /path/to/otel-go-demo
```

2. 启动所有服务：
```bash
docker-compose up --build
```

这将启动三个容器：
- **app**: Go Web 应用（端口 8080）
- **otel-collector**: OpenTelemetry Collector（端口 4317/4318）
- **otel-go-agent**: Go 自动插桩代理

3. 测试应用：
```bash
# 访问 Hello 端点
curl http://localhost:8080/hello

# 查看 Prometheus 指标
curl http://localhost:8080/metrics
```

### 本地开发运行

1. 安装依赖：
```bash
go mod download
```

2. 运行应用：
```bash
go run cmd/server/main.go cmd/server/JsonOutputsForLogger.go
```

3. 应用将在 `http://localhost:8080` 上启动

## 配置说明

### OpenTelemetry Collector 配置

`otel-collector-config.yaml` 配置了：

**接收器 (Receivers)**:
- OTLP gRPC (端口 4317)
- OTLP HTTP (端口 4318)
- Prometheus Simple - 从应用的 `/metrics` 端点抓取指标

**处理器 (Processors)**:
- `batch`: 批量处理遥测数据
- `memory_limiter`: 限制内存使用
- `filter/ottl`: 过滤 `/metrics` 端点自身的追踪和指标

**导出器 (Exporters)**:
- `debug`: 将遥测数据输出到控制台（详细模式）

### Docker Compose 配置

#### 应用容器 (app)
- 暴露端口 8080
- 依赖 otel-collector

#### OpenTelemetry Collector
- 使用官方 contrib 镜像
- 挂载配置文件
- 暴露 OTLP 接收端口

#### Go 自动插桩代理 (otel-go-agent)
- 使用 eBPF 技术自动插桩
- 需要特权模式和特定的 capabilities
- 自动追踪 HTTP 请求和运行时指标
- 环境变量配置：
  - `OTEL_EXPORTER_OTLP_ENDPOINT`: Collector 地址
  - `OTEL_SERVICE_NAME`: 服务名称
  - `OTEL_GO_AUTO_TARGET_EXE`: 目标可执行文件路径

## 日志格式

应用使用自定义的 JSON 日志处理器，输出格式如下：

```json
{
  "timestamp": "2024-11-30 12:00:00",
  "logging.googleapis.com/trace": "trace-id-here",
  "logging.googleapis.com/spanId": "span-id-here",
  "logging.googleapis.com/trace_sampled": true,
  "severity": "INFO",
  "message": "hello world!!!"
}
```

这种格式与 Google Cloud Logging 兼容，便于在云环境中使用。

## 开发指南

### 添加新的端点

在 `main.go` 中添加新的 HTTP 处理器：

```go
s.BindHandler("/your-endpoint", func(r *ghttp.Request) {
    g.Log().Info(r.Context(), "处理请求")
    r.Response.Write("响应内容")
})
```

### 自定义指标

使用 OpenTelemetry 指标 API：

```go
meter := gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
    Instrument:        "your-instrument",
    InstrumentVersion: "v1.0",
})

counter := meter.MustCounter("your.metric.name", gmetric.MetricOption{
    Help: "指标说明",
    Unit: "单位",
})

counter.Add(ctx, 1)
```

### 手动添加追踪（manual 分支）

在 `manual` 分支中，你可以手动创建和管理追踪 span：

#### 1. 初始化追踪提供者

在 `main` 函数中初始化 OTLP HTTP 导出器：

```go
import "github.com/gogf/gf/contrib/trace/otlphttp/v2"

const (
    serviceName = "your-service-name"
    endpoint    = "host.docker.internal:4318"  // OTel Collector 地址
    path        = "/v1/traces"                  // OTLP HTTP 路径
)

func main() {
    ctx := gctx.New()
    shutdown, err := otlphttp.Init(serviceName, endpoint, path)
    if err != nil {
        g.Log().Fatalf(ctx, "failed to initialize tracer: %v", err)
    }
    defer shutdown(ctx)
    // ...其他代码
}
```

#### 2. 创建自定义 Span

使用 OpenTelemetry API 创建 span：

```go
import (
    "go.opentelemetry.io/otel"
    "github.com/gogf/gf/v2/net/gtrace"
)

func yourBusinessLogic(ctx context.Context) {
    // 方式 1: 使用原生 OpenTelemetry API
    tracer := otel.GetTracerProvider().Tracer("your-tracer-name")
    ctx, span := tracer.Start(ctx, "your-operation-name")
    defer span.End()
    
    // 添加事件
    span.AddEvent("processing started")
    
    // 添加属性
    span.SetAttributes(
        attribute.String("user.id", "123"),
        attribute.Int("items.count", 10),
    )
    
    // 执行业务逻辑
    doSomething(ctx)
    
    // 方式 2: 使用 GoFrame 的 gtrace 包
    ctx, innerSpan := gtrace.NewSpan(ctx, "inner-operation")
    defer innerSpan.End()
    
    doAnotherThing(ctx)
}
```

#### 3. 嵌套 Span 示例

创建父子关系的 span 来追踪复杂的调用链：

```go
func processOrder(ctx context.Context, orderID string) {
    tracer := otel.GetTracerProvider().Tracer("order-service")
    
    // 父 span
    ctx, parentSpan := tracer.Start(ctx, "process-order")
    parentSpan.SetAttributes(attribute.String("order.id", orderID))
    defer parentSpan.End()
    
    // 子 span 1: 验证订单
    validateOrder(ctx, orderID)
    
    // 子 span 2: 处理支付
    processPayment(ctx, orderID)
    
    // 子 span 3: 发送通知
    sendNotification(ctx, orderID)
}

func validateOrder(ctx context.Context, orderID string) {
    ctx, span := gtrace.NewSpan(ctx, "validate-order")
    defer span.End()
    // 验证逻辑
}
```

#### 4. 记录错误

在 span 中记录错误信息：

```go
import "go.opentelemetry.io/otel/codes"

func riskyOperation(ctx context.Context) error {
    ctx, span := gtrace.NewSpan(ctx, "risky-operation")
    defer span.End()
    
    err := doSomethingRisky()
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }
    
    span.SetStatus(codes.Ok, "success")
    return nil
}
```

### 调试

查看 OpenTelemetry Collector 日志：
```bash
docker-compose logs -f otel-collector
```

查看应用日志：
```bash
docker-compose logs -f app
```

查看自动插桩代理日志：
```bash
docker-compose logs -f otel-go-agent
```

## 参考资源

- [GoFrame 官方文档](https://goframe.org/)
- [OpenTelemetry Go 文档](https://opentelemetry.io/docs/instrumentation/go/)
- [OpenTelemetry Go Auto-Instrumentation](https://github.com/open-telemetry/opentelemetry-go-instrumentation)
- [Prometheus 文档](https://prometheus.io/docs/)
- [OpenTelemetry Collector 文档](https://opentelemetry.io/docs/collector/)
