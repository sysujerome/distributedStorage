# # redis-server &

# function generate_data {
#     filename=$(pwd)"/data/benchmark_data"
#     if [ -e $filename ]
#     then
#         # touch $filename
#         > $filename
#         for i in $(seq 2 1 $1)
#         do
#             echo `num` >> $filename
#         done
#         echo -n `num` >> $filename
#     fi
# }

# function num {
#     let i=$RANDOM%3
#     opt=""
#     key=`echo $RANDOM | base64 | head -c 20`
#     value=$(echo $RANDOM | base64 | head -c 20)
#     if [ $i == 0 ] 
#     then
#         opt="get "$key
#     elif [ $i == 1 ] 
#     then
#         opt="set "$key" "$value
#     else
#         opt="del "$key
#     fi
#     echo $opt
# }


go build -o generate ./util/generateData/generateData.go
filename=$(pwd)"/data/benchmark_data"
if ! [ -e $filename ] 
then
    touch $filename
    ./generate $1
else
    length=(`wc -l $filename`)
    diff=`expr $1 - $length`
    if ! [ $diff == 1 ] 
    then
        echo "regenerate data file..."
        > $filename
        ./generate $1
    fi
fi

rm -rf ./kvserver
go build -o kvserver ./server/main.go
./kvserver --shard_node_name=shard_node_0 &
./kvserver --shard_node_name=shard_node_1 &
./kvserver --shard_node_name=shard_node_2 &
./kvserver --shard_node_name=shard_node_3 &
./kvserver --shard_node_name=shard_node_4 &
./kvserver --shard_node_name=shard_node_5 &
./kvserver --shard_node_name=shard_node_6 &
./kvserver --shard_node_name=shard_node_7 &
# ./kvserver --shard_node_name=shard_node_0 &

go build -o cli client/main.go
./cli 50050

ps -ef | grep kvserver | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
# ps -ef | grep storeInside | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
# ps -ef | grep storeInRedis | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
