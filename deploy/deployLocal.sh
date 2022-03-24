./deploy/shutdown.sh 

for idx in $(seq 0 15)
do
    # port=`expr 50050 + $idx`
    ./bin/server -shard_idx=$idx -conf=./conf_isee.json
done