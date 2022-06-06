# Kafka go demo

1. 对比

稳定性：sarama-cluster更高，相比而言concluent库会切到cgo，多了一层runtime的切换，不确定是否所有的异常都能完全转换到golang并处理

性能：concluent库更高，github上也有benchmark的结果。sarama使用人数比较多，但相对较为难用，性能较好
