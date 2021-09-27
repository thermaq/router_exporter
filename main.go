package main

import (
    "fmt"
    "net"
    "log"
    "strconv"
    "net/http"
    "github.com/jessevdk/go-flags"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket"
)

var (
    debug = true
    traffic = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "router_traffic",
        Help: "Traffic in bytes per IP",
    }, []string{"ip", "direction"})
)

func CountPacket(packet gopacket.Packet, ipnet *net.IPNet) {
    if tcpLayer := packet.Layer(layers.LayerTypeIPv4); tcpLayer != nil {
        iplayer, _ := tcpLayer.(*layers.IPv4)
        ipA := iplayer.SrcIP
        ipB := iplayer.DstIP
        length := float64(iplayer.Length)
        if debug {
            fmt.Printf("From %-15s To %-15s Length %-12d", ipA, ipB, int16(length))
        }
        ipAs := ipA.String()
        ipBs := ipB.String()
        if ipnet.Contains(ipA) {
            if ipnet.Contains(ipB) {
                // traffic inside LAN
                traffic.WithLabelValues(ipAs, "local").Add(length)
                traffic.WithLabelValues(ipBs, "local").Add(length)
                if debug {fmt.Printf("(local)\n")}
            } else {
                // trafic from inside to outside
                traffic.WithLabelValues(ipAs, "egress").Add(length)
                if debug {fmt.Printf("(egress)\n")}
            }
        } else if ipnet.Contains(ipB) {
            // traffic from outside to inside
            traffic.WithLabelValues(ipBs, "ingress").Add(length)
            if debug {fmt.Printf("(ingress)\n")}
        } else {
            // dont care
            if debug {fmt.Printf("(skip)\n")}
        }
    }
}

func RecordMetrics(interfac string, filter string, mask string) {
    _,ipnet,_ := net.ParseCIDR(mask)
    if handle, err := pcap.OpenLive(interfac, 1600, true, pcap.BlockForever); err != nil {
      panic(err)
    } else if err := handle.SetBPFFilter(filter); err != nil {
      panic(err)
    } else {
      packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
      for packet := range packetSource.Packets() {
        CountPacket(packet, ipnet)
      }
    }
}

type Options struct {
   Interface string `long:"interface" description:"Interface to listen on" default:"any"`
   Filter string `long:"filter" description:"pcap filter" default:""`
   Mask string `long:"mask" description:"ip range to gather data for" default:"0.0.0.0/0"`
   Debug bool `long:"debug" description:"print captured packets"`
   Port int `long:"port" description:"port to listen on with /metrics" default:"2112"`
   Address string `long:"address" description:"address to listen on with /metrics" default:"127.0.0.1"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    debug = opts.Debug
    go RecordMetrics(opts.Interface, opts.Filter, opts.Mask)

    address := opts.Address + ":" + strconv.Itoa(opts.Port)
    fmt.Println("Listening on ", address)
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(address, nil))
}

