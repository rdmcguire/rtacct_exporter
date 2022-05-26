# Linux Route Realm Exporter

This is a super simple exporter that does nothing but scrape and present packets and bytes from rtacct and provide counters for Prometheus at /metrics.

## Statistics

```
# HELP rtacct_bytes Counter of bytes from/to realm
# TYPE rtacct_bytes counter
rtacct_bytes{direction="from",realm="50West"} 380280
rtacct_bytes{direction="from",realm="mercury"} 5.367808e+06
rtacct_bytes{direction="from",realm="starlink"} 3.8832128e+07
rtacct_bytes{direction="to",realm="50West"} 4.4199936e+07
rtacct_bytes{direction="to",realm="mercury"} 1.062794e+06
rtacct_bytes{direction="to",realm="starlink"} 1.068864e+06
# HELP rtacct_pkts Counter of packets from/to realm
# TYPE rtacct_pkts counter
rtacct_pkts{direction="from",realm="50West"} 14824
rtacct_pkts{direction="from",realm="mercury"} 54940
rtacct_pkts{direction="from",realm="starlink"} 267369
rtacct_pkts{direction="to",realm="50West"} 322309
rtacct_pkts{direction="to",realm="mercury"} 5209
rtacct_pkts{direction="to",realm="starlink"} 9615
# HELP rtacct_stats_duration_us Time spend gathering statistics in microseconds
# TYPE rtacct_stats_duration_us histogram
rtacct_stats_duration_us_bucket{le="500"} 0
rtacct_stats_duration_us_bucket{le="1000"} 0
rtacct_stats_duration_us_bucket{le="1500"} 0
rtacct_stats_duration_us_bucket{le="2000"} 0
rtacct_stats_duration_us_bucket{le="2500"} 0
rtacct_stats_duration_us_bucket{le="3000"} 0
rtacct_stats_duration_us_bucket{le="3500"} 0
rtacct_stats_duration_us_bucket{le="4000"} 0
rtacct_stats_duration_us_bucket{le="4500"} 0
rtacct_stats_duration_us_bucket{le="5000"} 1
rtacct_stats_duration_us_bucket{le="5500"} 1
rtacct_stats_duration_us_bucket{le="6000"} 1
rtacct_stats_duration_us_bucket{le="6500"} 1
rtacct_stats_duration_us_bucket{le="7000"} 1
rtacct_stats_duration_us_bucket{le="7500"} 1
rtacct_stats_duration_us_bucket{le="8000"} 1
rtacct_stats_duration_us_bucket{le="8500"} 1
rtacct_stats_duration_us_bucket{le="9000"} 1
rtacct_stats_duration_us_bucket{le="9500"} 2
rtacct_stats_duration_us_bucket{le="10000"} 2
rtacct_stats_duration_us_bucket{le="+Inf"} 2
rtacct_stats_duration_us_sum 13883
rtacct_stats_duration_us_count 2500
```
