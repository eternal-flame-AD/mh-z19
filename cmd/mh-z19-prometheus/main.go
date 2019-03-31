package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"

	z19 "github.com/eternal-flame-AD/mh-z19"
)

var conn *serial.Port

var prefix string
var metricAddr string

type Measurement struct {
	Prefix        string
	Timestamp     int64
	Concentration uint16
}

var metrics Measurement
var metricsLock = sync.RWMutex{}
var metricTpl = template.Must(template.New("metric").Parse(strings.TrimSpace(`
# HELP {{.Prefix}}_co2 pressure in ppm
# TYPE {{.Prefix}}_co2 gauge
{{.Prefix}}_co2{device="mh-z19"} {{printf "%d" .Concentration}} {{.Timestamp}}
`)))

func init() {
	p := flag.String("p", "mh_z19", "metric prefix")
	m := flag.String("m", ":9100", "metrics address")
	dev := flag.String("d", "/dev/serial0", "device")
	flag.Parse()

	prefix = *p
	metricAddr = *m

	connConfig := z19.CreateSerialConfig()
	connConfig.ReadTimeout = 3 * time.Second
	connConfig.Name = *dev

	connSerial, err := serial.OpenPort(connConfig)
	if err != nil {
		panic(err)
	}
	conn = connSerial
}

func takeReading() {
	concentration, err := z19.TakeReading(conn)
	if err != nil {
		log.Println(err)
		return
	}
	now := time.Now()
	metricsLock.Lock()
	defer metricsLock.Unlock()
	metrics = Measurement{
		Concentration: concentration,
		Timestamp:     now.Unix() * 1000,
		Prefix:        prefix,
	}
}

func main() {

	takeReading()
	go func() {
		for range time.NewTicker(10 * time.Second).C {
			takeReading()
		}
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; version=0.0.4")
		metricsLock.RLock()
		defer metricsLock.RUnlock()
		metricTpl.Execute(w, metrics)
		w.Write([]byte{'\n'})
	})

	log.Println("promethus metrics running at " + metricAddr)
	log.Fatal(http.ListenAndServe(metricAddr, nil))

}
