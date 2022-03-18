current_date_time="`date +%Y/%m/%d' '%H:%M:%S`"
echo >> log
echo $current_date_time >> log

go build -o ./bin/server ./server/main.go
go build -o ./bin/client ./client/main.go
ps -ef | grep ./bin/server | grep -v grep | awk '{print $2}' | xargs kill -9
ps -ef | grep ./bin/client | grep -v grep | awk '{print $2}' | xargs kill -9
# ./bin/server --shard_idx=0 &>> log &
# ./bin/server --shard_idx=1 &>> log &
# ./bin/server --shard_idx=2 &>> log &
# ./bin/server --shard_idx=3 &>> log &
# ./bin/server --shard_idx=4 &>> log &
# ./bin/server --shard_idx=5 &>> log &
# ./bin/server --shard_idx=6 &>> log &
# ./bin/server --shard_idx=7 &>> log &

sleep 5s
count=$1
for i in $(seq 0  $[$count-1])
do
    ./bin/server --shard_idx=$i &>> log &
done
echo 


# ./bin/client --ip=192.168.1.128 --port=50050
echo "---------" >> log
