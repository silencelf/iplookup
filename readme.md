### fail2ban log analyze
### run the command to process raw log
grep -o -E '\d+\.\d+\.\d+\.\d+'  fail2ban.log | sort > ip.txt
### run the program
go run main.go
