package main

import "os"

var infxAddr = os.Getenv("INFLUXDB_ADDRESS")
var SQSURL = os.Getenv("AWS_SQS_URL")
var SQSRegion = os.Getenv("AWS_SQS_REGION")

func main() {
	mindSQL()
	mindNsq()
	mindInfluxDB()
	mindSQS()
}
