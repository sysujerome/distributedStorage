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


go build -o generate ./util/generateData.go
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


# redis-server &> /dev/null &> /dev/null &
# go build -o storeInRedis server/storeInRedis.go 
# go build -o storeInRemoteRedis server/storeInRemoteRedis.go 
# go build -o storeInside server/storeInside.go
# ./storeInRedis &
# ./storeInRemoteRedis  &
# ./storeInside &
# sleep 1s

# go build -o cli client/main.go
# ./cli 50050
# ./cli 50052
# ./cli 50051

# sleep 1s

ps -ef | grep redis-server | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
ps -ef | grep storeInside | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
ps -ef | grep storeInRedis | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
