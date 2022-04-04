chmod +x ./bin/storageInside
echo > log
./bin/storageInside --shard_idx=$1 --conf=conf.json & >> log