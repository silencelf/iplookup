# iplookup
### log analyze
First you should preprocess logs which contains ip address
example:
```
1.1.1.1
1.1.1.2
233.23.23.11
```
## Examples:
### run the command to process raw log:
grep -o -E '\d+\.\d+\.\d+\.\d+' logfile | uniq | sort > ip.txt
### installation via go install:
go install iplookup.go
### run the program:
go run iplookup.go 1.1.1.1
go run iplookup.go -f filename
### if it's installed:
iplookup 1.1.1.1
iplookup.go -f filename
### work together with unix pipelines:
grep -o -E '\d+\.\d+\.\d+\.\d+'  fail2ban.log | uniq | iplookup