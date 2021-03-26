// helloworld project main.go.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
        "strconv"
	"github.com/warthog618/gpio"
)


// Watches GPIO 4 (J8 7) and reports when it changes state.
func main() {
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()
	pin := gpio.NewPin(gpio.GPIO6)
	pin.Input()
	pin.PullUp()

  pinPWM1 := gpio.NewPin(gpio.GPIO18)
  pinPWM1.Output()

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

  signal, err := strconv.ParseFloat(os.Args[1], 64)
  signal = signal*100

	err = pin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
		if pin.Read() == true {
        time.Sleep((time.Duration(signal) * time.Microsecond))
        pinPWM1.High()
    }

    if pin.Read() == false {
      time.Sleep((time.Duration(signal) * time.Microsecond))
      pinPWM1.Low()
    }
	})

	if err != nil {
		panic(err)
	}

  defer pin.Unwatch()

	// In a real application the main thread would do something useful here.
	// But we'll just run for a minute then exit.
	fmt.Println("Watching Pin 4...")
	select {
  case <-time.After(time.Minute):
	case <-quit:
	}
}
