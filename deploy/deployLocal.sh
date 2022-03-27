./deploy/shutdown.sh 
#go build -o ./bin/server ./server/storeInside.go
echo > log
for idx in $(seq 0 15)
do
    # port=`expr 50050 + $idx`
    ./bin/storageInside -shard_idx=$idx -conf=./conf_isee.json &>> log &
done
