ip="10.182.183"
for sitenu in $(seq 1 16)
do
    ping -c 1 -w 1 ${ip}.${sitenu} &> /dev/null && result=0 || result=1
    if [ "${result}" == 0 ]; then
        echo "Server ${ip}.${sitenu} is UP."
    else
        echo "Server ${ip}.${sitenu} is Down."
    fi

done

# network="10.182.183"              # 先定义一个网域的前面部分！
# for sitenu in $(seq 1 16)       # seq 为 sequence（连续） 的缩写之意
# do
# 	# 下面的程序在取得 ping 的回传值是正确的还是失败的！
#         ping -c 1 -w 1 ${network}.${sitenu} &> /dev/null && result=0 || result=1
# 	# 开始显示结果是正确的启动 （UP） 还是错误的没有连通 （DOWN）
#         if [ "${result}" == 0 ]; then
#                 echo "Server ${network}.${sitenu} is UP."
#         else
#                 echo "Server ${network}.${sitenu} is DOWN."
#         fi
# done