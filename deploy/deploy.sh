
cluster_ip=(
    "10.182.183.5"
    "10.182.183.6"
    "10.182.183.7"
    "10.182.183.8"
    "10.182.183.9"
    "10.182.183.10"
    "10.182.183.11"
    "10.182.183.12"
    "10.182.183.13"
    "10.182.183.14"
    "10.182.183.15"
    "10.182.183.16"
)

cluster_name=(
    "thu@cn16611"
    "thu@cn16612"
    "thu@cn16613"
    "thu@cn16614"
    "thu@cn16615"
    "thu@cn16616"
    "thu@cn16617"
    "thu@cn16618"
    "thu@cn16619"
    "thu@cn16620"
    "thu@cn16621"
    "thu@cn16622"
)

path="/home/thu/distributedStorage"

./deploy/sync.sh

for idx in $(seq 0 11)
do

        ssh -t -t ${cluster_name[${idx}]} << EOF
        cd ${path}
        chmod +x ./deploy/shutdown.sh
        chmod +x ./deploy/start.sh
        ./deploy/shutdown.sh
        ./deploy/start.sh ${idx}

EOF

done

