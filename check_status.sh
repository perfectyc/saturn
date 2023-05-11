go run /root/saturn/check_active.go > status.txt

cat status.txt | grep -w   "active"  >  status_active.txt
cat status.txt | grep -w   "inactive" > status_inactive.txt
cat status.txt | grep -w   "down"  > status_down.txt
cat status.txt | grep -w   "found"  > status_no_found.txt
cat status_active.txt && cat status_inactive.txt && cat status_down.txt && cat status_no_found.txt
