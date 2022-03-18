#! /usr/bin/bash
# set -x # 把执行的命令在终端显示出来

YCSB=/home/pjl/application/ycsb-0.17.0/bin/ycsb  #shell中变量赋值时变量和等号之间不能有空格
WL_PATH=/home/pjl/lab/distributedStorage/benchmark/workloads_20w
mkdir -p ${WL_PATH} # -p 文件夹不存在则创建

#1M-16M
WL_TYPES=(a b c d) #数组
# WL_TYPES=(a)
#WL_TYPE=a
WORKLOAD_COUNT=200000

# for WORKLOAD_DIST in ${WORKLOAD_DISTS[*]}
# do

for WL_TYPE in ${WL_TYPES[*]} # [*]表示数组的所有元素
do

WORKLOAD=workload${WL_TYPE} # 拼接

TEMPWORKLOADFILENAME=${WL_PATH}/${WORKLOAD}-${WORKLOAD_COUNT} # 拼接;这个是参数集合

if [ ${WORKLOAD} == "workloada" ]
then

# generate workload config

# Yahoo! Cloud System Benchmark
# Workload A: Update heavy workload
#   Application example: Session store recording recent actions
#
#   Read/update ratio: 50/50
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: zipfian

# 多行导入
cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=0.5
updateproportion=0.5
scanproportion=0
insertproportion=0
hotspotdatafraction=0.1
hotspotopnfraction=0.9

requestdistribution=hotspot
EOF


cat ${TEMPWORKLOADFILENAME}


elif [ ${WORKLOAD} == "workloadb" ]
then

# generate workload config
# Yahoo! Cloud System Benchmark
# Workload B: Read mostly workload
#   Application example: photo tagging; add a tag is an update, but most operations are to read tags
#
#   Read/update ratio: 95/5
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: zipfian

cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=0.95
updateproportion=0.05
scanproportion=0
insertproportion=0
hotspotdatafraction=0.1
hotspotopnfraction=0.9

requestdistribution=hotspot
EOF

cat ${TEMPWORKLOADFILENAME}

elif [ ${WORKLOAD} == "workloadc" ]
then

# generate workload config
# Yahoo! Cloud System Benchmark
# Workload C: Read only
#   Application example: user profile cache, where profiles are constructed elsewhere (e.g., Hadoop)
#
#   Read/update ratio: 100/0
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: zipfian

cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=1
updateproportion=0
scanproportion=0
insertproportion=0
hotspotdatafraction=0.1
hotspotopnfraction=0.9

requestdistribution=hotspot
EOF

cat ${TEMPWORKLOADFILENAME}


elif [ ${WORKLOAD} == "workloadd" ]
then

# generate workload config
# Yahoo! Cloud System Benchmark
# Workload D: Read latest workload
#   Application example: user status updates; people want to read the latest
#
#   Read/update/insert ratio: 95/0/5
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: latest

# The insert order for this is hashed, not ordered. The "latest" items may be
# scattered around the keyspace if they are keyed by userid.timestamp. A workload
# which orders items purely by time, and demands the latest, is very different than
# workload here (which we believe is more typical of how people build systems.)

cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=0.95
updateproportion=0
scanproportion=0
insertproportion=0.05
hotspotdatafraction=0.1
hotspotopnfraction=0.9

requestdistribution=hotspot
EOF

cat ${TEMPWORKLOADFILENAME}

elif [ ${WORKLOAD} == "workloade" ]
then

# generate workload config
# Yahoo! Cloud System Benchmark
# Workload E: Short ranges
#   Application example: threaded conversations, where each scan is for the posts in a given thread (assumed to be clustered by thread id)
#
#   Scan/insert ratio: 95/5
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: zipfian

# The insert order is hashed, not ordered. Although the scans are ordered, it does not necessarily
# follow that the data is inserted in order. For example, posts for thread 342 may not be inserted contiguously, but
# instead interspersed with posts from lots of other threads. The way the YCSB client works is that it will pick a start
# key, and then request a number of records; this works fine even for hashed insertion.

cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=0
updateproportion=0
scanproportion=0.95
insertproportion=0.05

requestdistribution=zipfian

maxscanlength=100

scanlengthdistribution=uniform
EOF

cat ${TEMPWORKLOADFILENAME}

elif [ ${WORKLOAD} == "workloadf" ]
then

# generate workload config
# Yahoo! Cloud System Benchmark
# Workload F: Read-modify-write workload
#   Application example: user database, where user records are read and modified by the user or to record user activity.
#
#   Read/read-modify-write ratio: 50/50
#   Default data size: 1 KB records (10 fields, 100 bytes each, plus key)
#   Request distribution: zipfian


cat <<EOF > ${TEMPWORKLOADFILENAME}
fieldcount=1
fieldlength=1
recordcount=${WORKLOAD_COUNT}
operationcount=${WORKLOAD_COUNT}
workload=site.ycsb.workloads.CoreWorkload

readallfields=true

readproportion=0.5
updateproportion=0
scanproportion=0
insertproportion=0
readmodifywriteproportion=0.5

requestdistribution=zipfian
EOF

cat ${TEMPWORKLOADFILENAME}



else

echo ${WORKLOAD} "is not implemented !"
exit

fi 


LOADLOGFILE=${WL_PATH}/${WORKLOAD}-load-${WORKLOAD_COUNT}.log
# ${YCSB} load basic -P ${TEMPWORKLOADFILENAME}
${YCSB} load basic -P ${TEMPWORKLOADFILENAME} > ${LOADLOGFILE} 2>&1
./log_to_workload.py ${LOADLOGFILE}
rm ${LOADLOGFILE}

RUNLOGFILE=${WL_PATH}/${WORKLOAD}-run-${WORKLOAD_COUNT}.log
${YCSB} run basic -P ${TEMPWORKLOADFILENAME} > ${RUNLOGFILE} 2>&1
./log_to_workload.py ${RUNLOGFILE}
rm ${RUNLOGFILE}

# DELETE Intermediate
rm ${TEMPWORKLOADFILENAME}
done
