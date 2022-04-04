cluster_name=(
    "thu@cn16607"
    "thu@cn16608"
    "thu@cn16609"
    "thu@cn16610"
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

path="/home/thu/lab/deploy"
ssh -t -t ${cluster_name[$1]} << EOF
    echo "log>>>>>>>>>>>>>>>>>>>>>"
    cat $path/log
    exit
EOF