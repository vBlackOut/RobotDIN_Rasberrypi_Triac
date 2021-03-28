package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
        "bufio"
        "strconv"
	"github.com/warthog618/gpio"
        "strings"
)


func getInt(key string) uint32 {
    //s := os.Getenv(key)
    if signals, err := strconv.ParseFloat(key, 32); err == nil {

      if signals <= 100 {
        signals = 10000-(signals*100)
      } else if (signals <= 1000) && (signals > 100) {
        signals = 1000-(signals*10)
      } else if (signals <= 10000) && (signals > 1000) {
        signals = (signals)
      }

      return uint32(signals)
    }
    return 0
}

func getenvInt(key string) uint32 {
    s := os.Getenv(key)
    if signals, err := strconv.ParseFloat(s, 32); err == nil {

      if signals <= 100 {
        signals = 10000-(signals*100)
      } else if (signals <= 1000) && (signals > 100) {
        signals = 1000-(signals*10)
      } else if (signals <= 10000) && (signals > 1000) {
        signals = (signals)
      }

      return uint32(signals)
    }
    return 0
}

func pwm(pins int, pins2 int) (string, error) {

  err := gpio.Open()
  if err != nil {
    panic(err)
  }
  defer gpio.Close()
  pin := gpio.NewPin(gpio.GPIO6)
  pin.Input()
  pin.PullUp()

  pinPWM1 := gpio.NewPin(pins)
  pinPWM1.Output()

  pinPWM2 := gpio.NewPin(pins2)
  pinPWM2.Output()

  // capture exit signals to ensure resources are released on exit.
  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(quit)

  percent1 := uint32(0)
  percent2 := uint32(0)

  fmt.Println("init pwm process...")

  err = pin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
    if pin.Read() == true {
        if (percent1 > 0) {
          time.Sleep((time.Duration(percent1) * time.Microsecond))
          pinPWM1.High()
        }
        if (percent2 > 0) {
          time.Sleep((time.Duration(percent2) * time.Microsecond))
          pinPWM2.High()
        }
    }

    if pin.Read() == false {
      if (percent1 > 0) {
        time.Sleep((time.Duration(percent1) * time.Microsecond))
        pinPWM1.Low()
      }
      if (percent2 > 0) {
        time.Sleep((time.Duration(percent2) * time.Microsecond))
        pinPWM2.Low()
      }
    }
  })

  if err != nil {
    panic(err)
  }

  defer pin.Unwatch()

  // In a real application the main thread would do something useful here.
  // But we'll just run for a minute then exit.
  fmt.Println("please wait command usage :\n cmd : pwm[1-2] percent[0-100%]\n cmd : [reset] for disable all comand \n cmd : [stop] for quit program ")
  input := bufio.NewScanner(os.Stdin)
  for input.Scan() {
    if input.Text() == "stop" {
        break
    } else {

      if strings.Contains(input.Text(), "pwm1") {
        detectpercent := strings.Split(input.Text(), " ")[1]
        percent1 = getInt(detectpercent)
        if percent1 == 10000 {
          percent1 = 0
        } else if percent1 == 0 {
          percent1 = 1
        }
        //fmt.Println("debug:", percent1)
      }

      if strings.Contains(input.Text(), "pwm2") {
        detectpercent := strings.Split(input.Text(), " ")[1]
        percent2 = getInt(detectpercent)
        if percent2 == 10000 {
          percent2 = 0
        } else if percent2 == 0 {
          percent2 = 1
        }
        //fmt.Println("debug:", percent2)
      }

      if strings.Contains(input.Text(), "reset") {
        percent1 = uint32(0)
        percent2 = uint32(0)
      }
    }
  }
  // select {
  // case <-time.After(time.Duration(timeout) * time.Second):
  // case <-quit:
  // }
  return "", nil
}

// Watches GPIO 4 (J8 7) and reports when it changes state.
func main() {
    pwm(gpio.GPIO18, gpio.GPIO19)
}
