package main

import (
    "log"
    "time"
    influx "github.com/influxdata/influxdb/client/v2"
)

var (
    MetricsClient influx.Client
)

func MetricsSend(l *LogEntry) {
    var (
	err error
    )
    
    bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{})
    
    tags := map[string]string{
	"Server":	l.Server,
	"Scheme":	l.Scheme,
	"Method":	l.Method,
	"Hostname":	l.Hostname,
	"Status":	l.Status,
	"Protocol":	l.Protocol,
	"Completed":	l.Completed,
	"Extension":	l.Extension,
    }
    
    fields := map[string]interface{}{
	"ClientIP":		l.ClientIP,
	"URI":			l.URI,
	"Duration":		l.Duration,
	"BytesSent":		int(l.BytesSent),
	"TCPInfo_RTT":		int(l.TCPInfo_RTT),
	"TCPInfo_RTTVAR":	int(l.TCPInfo_RTTVAR),
	"TCPInfo_SND_CWND":	int(l.TCPInfo_SND_CWND),
	"TCPInfo_RCV_Space":	int(l.TCPInfo_RCV_Space),
    }
    
    pt, err := influx.NewPoint(measurement, tags, fields, time.Now())
    if err != nil {
	log.Println("Unable to create datapoint: " + err.Error())
	return
    }
    
    bp.AddPoint(pt)
    
    if err = MetricsClient.Write(bp); err != nil {
	log.Println("Unable to write: " + err.Error())
    }
}

func MetricsInit() {
    MetricsClient, _ = influx.NewUDPClient(influx.UDPConfig{Addr: hostport_influxdb})
}
