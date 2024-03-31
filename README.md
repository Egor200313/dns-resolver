# Simple dns-resolver

UDP server acting like dns resolver answering with records from given file only

## How to run

`go run main.go <records filename>` where `filename` is a file with dns records like records.txt in the same format

## To test server run

`dig @localhost -p 5003 <args>`

Example:

`dig @localhost -p 5003 example.com AAAA`

## Generate huge file with random A records

`go run generator/main.go` It will produce file `random_records.txt` in correct format with 10001 random A records