# Process Exporter

## 监控指标
`CPU`, `Memory`, `IO`, `NetConnection`, `FD`, `Thread`, `ContextSwitch`

## 通用标签
* `pid`: 进程pid
* `cmdline`: 进程名称

## CPU Metrics
### process_cpu_seconds_total
* 进程cpu使用时间, 除通用标签外, 有个`mode`标签, 取值`iowait`,`user`,`system`
* 查询cpu使用率
  * 按cpu类型: `sum(rate(process_cpu_seconds_total[1m])) by (pid, cmdline, mode)`
  * 总cpu:  `sum(rate(process_cpu_seconds_total[1m])) by (pid, cmdline)`
  

## Memory Metrics
### process_memory_rss_bytes
* 进程RSS内存使用量, 通用标签
* 查询表达式: `process_memory_rss_bytes`

## Network Metrics
### process_network_tcp_connections
* 进程网络TCP连接数, 通用标签
* 查询表达式: `process_network_tcp_connections`

### process_network_tcp_connections_status
* 进程网络TCP连接状态, 除通用标签外, 有个`status`标签
* 查询表达式: `process_network_tcp_connections_status`

### process_network_receive_bytes_total
* 进程网络接收流量, 通用标签
* 查询表达式: `rate(process_network_receive_bytes_total[1m])`

### process_network_transmit_bytes_total
* 进程网络发送流量, 通用标签
* 查询表达式: `rate(process_network_transmit_bytes_total[1m])`


## IO Metrics
### process_io_read_count_total
* 进程读io总量, 通用标签
* read iops查询表达式: `rate(process_io_read_count_total[1m])`

### process_io_write_count_total
* 进程写io总量, 通用标签
* write iops查询表达式: `rate(process_io_write_count_total[1m])`

### process_io_read_bytes_total
* 进程读bytes总量, 通用标签
* rBytes/s查询表达式: `rate(process_io_read_bytes_total[1m])`

### process_io_write_count_total
* 进程写bytes总量, 通用标签
* wBytes/s查询表达式: `rate(process_io_write_bytes_total[1m])`

## Common Metrics
### process_ctx_switch_voluntary_count_total
* 进程主动上下文切换, 通用标签
* 每秒主动上下文切换查询表达式: `rate(process_ctx_switch_voluntary_count_total[1m])`

### process_ctx_switch_involuntary_count_total
* 进程被动上下文切换, 通用标签
* 每秒被动上下文切换查询表达式: `rate(process_ctx_switch_involuntary_count_total[1m])`

### process_fds_count
* 进程文件描述符数量, 通用标签
* 查询表达式 `process_fds_count`

### process_threads_count
* 进程线程数量, 通用标签
* 查询表达式 `process_threads_count`

### process_total_count
* 总进程数量, 无标签
* 查询表达式 `process_total_count`

## 启动应用
`./process-exporter -listen-port 10921 -cmdline-include name1,name2,name3 -cmdline-exclude name1,name2,name3`

### 参数说明
* `listen-port`: 监听端口, 默认10921
* `cmdline-include`:进程名称匹配关键字, 不设置则收集所有进程, 设置的话**只会**收集进程名称含有关键字的进程, 多个关键字以`,`隔开, 推荐设置,全量收集比较耗资源
* `cmdline-exclude`:进程名称匹配关键字, 不设置则收集所有进程, 设置的话**不会**收集进程名称含有关键字的进程, 多个关键字以`,`隔开, 推荐设置,全量收集比较耗资源
* `cmdline-include`和`cmdline-exclude`同时设置的话, `cmdline-exclude`优先生效

## Grafana
可直接导入`Process Exporter-grafana.json`