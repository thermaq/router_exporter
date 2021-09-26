# Building

`go build`

# Usage

```
Usage:
  router_exporter [OPTIONS]

Application Options:
      --interface= Interface to listen on (default: any)
      --filter=    pcap filter
      --mask=      ip range to gather data for (default: 0.0.0.0/0)
      --debug      print captured packets
      --port=      port to listen on with /metrics (default: 2112)
      --address=   address to listen on with /metrics (default: 127.0.0.1)

Help Options:
  -h, --help       Show this help message

2021/09/26 11:43:28 Usage:
  router_exporter [OPTIONS]

Application Options:
      --interface= Interface to listen on (default: any)
      --filter=    pcap filter
      --mask=      ip range to gather data for (default: 0.0.0.0/0)
      --debug      print captured packets
      --port=      port to listen on with /metrics (default: 2112)
      --address=   address to listen on with /metrics (default: 127.0.0.1)

Help Options:
  -h, --help       Show this help message
```

# Example output

```
# TYPE router_traffic counter
router_traffic{direction="egress",ip="10.0.0.7"} 16996
router_traffic{direction="ingress",ip="10.0.0.7"} 3.004459e+06
```
