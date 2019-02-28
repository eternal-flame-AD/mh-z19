# mh-z19

Reading [mh-z19](https://www.winsen-sensor.com/products/ndir-co2-sensor/mh-z19.html) sensors with wiringpi and golang on raspberry pi.


## Example

```golang
package main

import (
	"fmt"

	z19 "github.com/eternal-flame-AD/mh-z19"
	"github.com/tarm/serial"
)

func main() {
	connConfig := z19.CreateSerialConfig()
	connConfig.Name = "/dev/serial0"

	port, err := serial.OpenPort(connConfig)
	if err != nil {
		panic(err)
	}
	defer port.Close()
	concentration, err := z19.TakeReading(port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("co2=%d ppm\n", concentration)
}

```