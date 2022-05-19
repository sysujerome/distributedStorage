# distributedStorage


本章提出了一种分布式的线性哈希算法，并基于gRPC框架实现了基于分布式的线性哈希的索引系统，使得索引系统可以进行自动按需扩容。针对分布式的元数据同步问题，本章提出了一种元数据同步策略，即集群内部的节点的元数据允许彼此间不一致，节点在被访问时确定当前的next和level值是否一致来做更新，客户端和索引系统的元数据不做主动更新，等索引系统在一轮分裂过后再被动更新元数据，避免了元数据不一致导致数据紊乱的问题。针对线性哈希的分裂触发机制会引起分裂操作爆发性的问题，本系统提出了动态扩缩阈值的方法，避免了集群在一瞬间发生大量分裂操作。为了减少节点分裂加锁过大地影响系统性能，本系统采取了把锁细化到指针的移动的方法。针对集群节点分裂时节点数据量较大的问题，本系统采用了针对每个槽建立一个线程来处理数据，将分裂节点耗费时间缩短了95\%。
