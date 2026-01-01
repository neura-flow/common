## worker

任务管理模块，提供多种类型的任务处理方式

#### Job / ContextJob / DelayJob / CronJob
定义了 Job、ContextJob、DelayJob、CronJob 等接口，提供实时任务、延时任务、定时任务等功能
- [代码](job.go)
- [示例](job_test.go)

#### Worker / CronWorker
定义了 Worker、CronWorker，用于执行各种任务
- [代码](worker.go)
- [示例](worker_test.go)