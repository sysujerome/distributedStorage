./deploy/shutdown.sh 
#go build -o ./bin/server ./server/storeInside.go
rm -rf ./overflow
echo > log
for idx in $(seq 0 15)
do
    # port=`expr 50050 + $idx`
    ./bin/server -shard_idx=$idx -conf=./conf_isee.json &>> log &
done
