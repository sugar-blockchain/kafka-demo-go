# Kafka go demo

此仓库包含两种Kafka客户端的实现，分别是sarama-cluster和confluent

## 对比

稳定性：sarama-cluster更高，相比而言confluent库会切到cgo，多了一层runtime的切换，不确定是否所有的异常都能完全转换到golang并处理

性能：confluent库更高，github上也有benchmark的结果。sarama使用人数比较多，但相对较为难用，性能较好

## 使用说明

[sarama-cluster](./client/my_sarama/README.md)

[confluent](./client/my_confluent/README.md)
