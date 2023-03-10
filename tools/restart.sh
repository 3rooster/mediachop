ps aux|grep mediachop|grep -v grep | awk '{print $2}'|xargs kill -9
GOTRACEBACK=crash nohup ./mediachop >> ../log/start.log   2>&1 &