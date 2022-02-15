# redis-server &

function generate_data {
    filename=$(pwd)"/data/benchmark_data"
    if [ -e $filename ]
    then
        # touch $filename
        > $filename
        for i in $(seq 2 1 $1)
        do
            echo `num` >> $filename
        done
        echo -n `num` >> $filename
    fi
}

function num {
    let i=$RANDOM%3
    opt=""
    key=`echo $RANDOM | base64 | head -c 20`
    value=$(echo $RANDOM | base64 | head -c 20)
    if [ $i == 0 ] 
    then
        opt="get "$key
    elif [ $i == 1 ] 
    then
        opt="set "$key" "$value
    else
        opt="del "$key
    fi
    echo $opt
}


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


redis-server &> /dev/null &
go run server/main.go &> /dev/null &
go run server/storeInside.go   &> /dev/null &
go build -o cli client/main.go
./cli 50050
./cli 50051

sleep 1s

ps -ef | grep redis-server | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
ps -ef | grep server/storeInside.go | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
ps -ef | grep server/main.go | grep -v grep | awk '{print $2}' | xargs kill -9 &> /dev/null &
