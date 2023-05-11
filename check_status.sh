go run /root/saturn/check_status.go > /root/saturn/status.txt

cat /root/saturn/status.txt | grep -w   "active"  >  /root/saturn/status_active.txt
cat /root/saturn/status.txt | grep -w   "inactive" > /root/saturn/status_inactive.txt
cat /root/saturn/status.txt | grep -w   "down"  > /root/saturn/status_down.txt
cat /root/saturn/status.txt | grep -w   "found"  > /root/saturn/status_no_found.txt
cat /root/saturn/status_active.txt && cat /root/saturn/status_inactive.txt && cat /root/saturn/status_down.txt && cat /root/saturn/status_no_found.txt
