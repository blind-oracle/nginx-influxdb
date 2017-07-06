# nginx-influxdb

This small utility receives syslog messages (access log) from nginx over UDP, parses them and sends to InfluxDB (or telegraf) over UDP.
You can use telegraf's logparser plugin instead, but this is more handy i think.

You should use the following nginx configuration:
```
map $status $metrics {
    ~^[23]  1;
    default 0;
}

log_format collector '$msec|$remote_addr|$scheme|$request_method|$host|$request_uri|$status|$request_time|$bytes_sent|$server_protocol|$completed|$tcpinfo_rtt|$tcpinfo_rttvar|$tcpinfo_snd_cwnd|$tcpinfo_rcv_space';
access_log syslog:server=1.1.1.1:514,tag=nginx collector if=$metrics;
```
