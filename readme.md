### fail2ban log analyze
First you should preprocess logs which contains ip address
example:
```
1.1.1.1
1.1.1.2
233.23.23.11
```
### run the command to process raw log
grep -o -E '\d+\.\d+\.\d+\.\d+'  fail2ban.log | sort > ip.txt
### run the program
go run ipsearch.go
### installation
go install ipsearch.go
