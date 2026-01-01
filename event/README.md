## event

事件驱动模块，提供各种事件处理组件

#### Event

定义了事件和事件处理函数接口
- [代码](event.go)
- [示例](event_test.go)

#### Subscriber

定义了订阅者接口，用于根据各种条件订阅事件
- [代码](subscriber.go)
- [示例](subscriber_test.go)

#### Reporter
定义了报告器接口，用于将事件的相关信息报告给其他模块/系统。

内置了 MetricsReporter、LoggerReporter、AlarmReporter
- [代码](reporter.go)
- [示例](reporter_test.go)

#### Interceptor
定义了拦截器，用于在处理事件前/后做一些通用的事情，比如过滤、输出日志、监控指标等
- [代码](interceptor.go)
- [示例](interceptor_test.go)

#### Source 和 Sink

定义了事件源(Source)和事件池(Sink)，事件源用于拉取事件，事件池用于推送事件
- [代码](sink_source.go)
- [示例](sink_source_test.go)

#### Selector

定义了事件选择器，用于从多个事件源中接收事件
- [代码](selector.go)
- [示例](selector_test.go)

#### Consumer
实现了简单的事件消费者，组合了 Subscriber 来支持按条件从事件源消费事件
- [代码](consumer.go)
- [示例](consumer_test.go)
