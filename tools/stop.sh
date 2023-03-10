
ps aux|grep mediachop|grep -v grep | awk '{print $2}'|xargs kill -9