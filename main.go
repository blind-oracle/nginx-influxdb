package main

import (
    "flag"
    "fmt"
    "time"
    "log"
    "strings"
    "strconv"
    
    syslog "gopkg.in/mcuadros/go-syslog.v2"
)

type LogEntry struct {
    TimeStamp		time.Time
    
    // Tags
    Server		string
    Scheme		string
    Method		string
    Hostname		string
    Status		string
    Protocol		string
    Completed		string
    Extension		string
    
    // Fields
    ClientIP		string
    URI			string
    Duration		float64
    BytesSent		uint64
    TCPInfo_RTT		uint64
    TCPInfo_RTTVAR	uint64
    TCPInfo_SND_CWND	uint64
    TCPInfo_RCV_Space	uint64
}

var (
    hostport_listen string
    hostport_influxdb string
    measurement	string
    debug bool
)

func main() {
    flag.StringVar(&hostport_listen, "listen", "0.0.0.0:514", "ip:port to listen for syslog messages")
    flag.StringVar(&hostport_influxdb, "influxdb", "", "ip:port of InfluxDB server to send metrics to")
    flag.StringVar(&measurement, "measurement", "", "InfluxDB measurement")
    flag.BoolVar(&debug, "debug", false, "Enable debug")
    
    flag.Parse()
    
    if hostport_influxdb == "" {
	log.Fatal("You must specify 'influxdb'")
    }
    
    if measurement == "" {
	log.Fatal("You must specify 'measurement'")
    }
    
    MetricsInit()
    
    channel := make(syslog.LogPartsChannel)
    handler := syslog.NewChannelHandler(channel)
    
    server := syslog.NewServer()
    server.SetFormat(syslog.RFC3164)
    server.SetHandler(handler)
    server.ListenUDP(hostport_listen)
    server.Boot()
    
    go func(channel syslog.LogPartsChannel) {
	for msg := range channel {
	    a := strings.Split(msg["content"].(string), "|")
	    
	    if len(a) != 15 {
		log.Println("Wrong number of fields")
		continue
	    }
	    
	    l := &LogEntry{
		Extension: "Unknown",
	    }
	    
	    l.Server = msg["hostname"].(string)
	    
	    l.ClientIP = a[1]
	    l.Scheme = a[2]
	    l.Method = a[3]
	    l.Hostname = a[4]
	    
	    b := strings.Split(a[5], "?")
	    l.URI = b[0]
	    
	    if !strings.HasSuffix(l.URI, "/") {
		c := strings.Split(l.URI, "/")
		if d := strings.Split(c[len(c)-1], "."); len(d) > 1 {
		    l.Extension = strings.ToLower(d[len(d)-1])
		}
	    }
	    
	    l.Status = a[6]
	    l.Duration, _ = strconv.ParseFloat(a[7], 64)
	    l.BytesSent, _ = strconv.ParseUint(a[8], 10, 64)
	    l.Protocol = a[9]
	    l.Completed = a[10]
	    l.TCPInfo_RTT, _ = strconv.ParseUint(a[11], 10, 64)
	    l.TCPInfo_RTTVAR, _ = strconv.ParseUint(a[12], 10, 64)
	    l.TCPInfo_SND_CWND, _ = strconv.ParseUint(a[13], 10, 64)
	    l.TCPInfo_RCV_Space, _ = strconv.ParseUint(a[14], 10, 64)
	    
	    if debug {
		fmt.Printf("%+v\n", l)
	    }
	    
	    MetricsSend(l)
	}
    }(channel)
    
    server.Wait()
}
