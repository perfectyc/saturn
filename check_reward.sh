go run /root/saturn/check_reward.go > /root/saturn/reward.txt
cat /root/saturn/reward.txt | grep -i -C 2 "Bogota" | grep total | awk '{print $NF}' > /root/saturn/bogota.txt
cat /root/saturn/reward.txt | grep -i -C 2 "lima" | grep total | awk '{print $NF}' >  /root/saturn/lima.txt
cat /root/saturn/reward.txt | grep -i -C 2 "sao paulo" | grep total | awk '{print $NF}' >  /root/saturn/paulo.txt
cat /root/saturn/reward.txt | grep -i -C 2 "santiago" | grep total | awk '{print $NF}' > /root/saturn/santiago.txt
cat /root/saturn/reward.txt | grep -i -C 2 "indon" | grep total | awk '{print $NF}' > /root/saturn/indon.txt
cat /root/saturn/reward.txt | grep -i -C 2 "hland" | grep total | awk '{print $NF}' > /root/saturn/hland.txt
total=0
for file in /root/saturn/bogota.txt /root/saturn/lima.txt /root/saturn/paulo.txt /root/saturn/santiago.txt /root/saturn/indon.txt /root/saturn/hland.txt; do
    result=$(awk '{sum+=$1} END {print sum}' "$file")
    echo "${file%%.*}过去24小时收益：$result"
    total=$(awk "BEGIN {print $total + $result}")
done
echo "所有地区过去24小时总收益：$total"
