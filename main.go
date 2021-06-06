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
	"sync"
)

func getInt(key string) uint32 {
    //s := os.Getenv(key)
    if signals, err := strconv.ParseFloat(key, 32); err == nil {

      if signals < 100 {
        signals = 9000-(signals*80)
      } else if (signals < 9000) && (signals > 1000) {
        signals = (signals)
      }

      if signals <= 1000 {
        signals = 1000
      }

      // if signals >= 9000 {
      // 	signals = 8700
      // }
      return uint32(signals)
    }
    return 0
}

func pwm2defer(wg *sync.WaitGroup, percent2 uint32) {
	if os.Getenv("pwmSTOP") == "True" {
		gpio.Close()
		os.Exit(0)
	}
	var pinPWM2 = gpio.NewPin(gpio.GPIO20)
	pinPWM2.Output()

	defer wg.Done()
	if ((percent2 >= 1) && (percent2 < 10000))  {
		time.Sleep((time.Duration(percent2) * time.Microsecond))
		pinPWM2.High()
		time.Sleep((time.Duration(5) * time.Microsecond))
		pinPWM2.Low()
	}
}

func pwm1defer(wg *sync.WaitGroup, percent1 uint32) {
	if os.Getenv("pwmSTOP") == "True" {
			gpio.Close()
			os.Exit(0)
	}
	var pinPWM1 = gpio.NewPin(gpio.GPIO16)
	pinPWM1.Output()
	defer wg.Done()
	if ((percent1 >= 1) && (percent1 < 10000))  {
		time.Sleep((time.Duration(percent1) * time.Microsecond))
		pinPWM1.High()
		time.Sleep((time.Duration(5) * time.Microsecond))
		pinPWM1.Low()
	}

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

  // capture exit signals to ensure resources are released on exit.
  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(quit)

  percent1 := uint32(0)
  percent2 := uint32(0)

  fmt.Println("init pwm process...")

  err = pin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
    //if pin.Read() == true {
	var wg sync.WaitGroup
	wg.Add(2)
	go pwm1defer(&wg, percent1)
	go pwm2defer(&wg, percent2)
	wg.Wait()

        // if ((percent1 >= 1) && (percent1 < 10000))  {
        //   time.Sleep((time.Duration(percent1) * time.Microsecond))
        //   pinPWM1.High()
        //   time.Sleep((time.Duration(5) * time.Microsecond))
        //   pinPWM1.Low()
        // }
				//
        // if ((percent2 >= 1) && (percent2 < 10000)) {
        //   time.Sleep((time.Duration(percent2) * time.Microsecond))
        //   pinPWM2.High()
        //   time.Sleep((time.Duration(5) * time.Microsecond))
        //   pinPWM2.Low()
        // }
  //  }

    // if pin.Read() == false {
    //   if ((percent1 >= 1) && (percent1 < 10000)) {
    //     time.Sleep((time.Duration(percent1) * time.Microsecond))
    //     pinPWM1.High()
    //     time.Sleep((time.Duration(5) * time.Nanosecond))
    //     pinPWM1.Low()
    //   }
    //   if ((percent2 >= 1) && (percent2 < 10000)) {
    //     time.Sleep((time.Duration(percent2) * time.Microsecond))
    //     pinPWM2.High()
    //     time.Sleep((time.Duration(5) * time.Nanosecond))
    //     pinPWM2.Low()
    //   }
    // }

		// if ((percent1 == 10000) || (percent1 == 0)) {
    //   pinPWM1.Low()
		// }
    //
		// if ((percent2 == 10000) || (percent2 == 0)) {
    //   pinPWM2.Low()
		// }

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

    fmt.Println(os.Getenv("pwmSTOP"))

    if input.Text() == "stop" || os.Getenv("pwmSTOP") == "True" {
	percent1 = uint32(0)
	percent2 = uint32(0)
	gpio.Close()
	break
    } else {

    if strings.Contains(input.Text(), "pwm1") {
        detectpercent1 := strings.Split(input.Text(), " ")[1]
        percent1 = getInt(detectpercent1)
	if ((percent1 == 8999) || (percent1 == 9000)) {
          percent1 = 0
        }
	// if percent1 == 0 {
        //   percent1 = 0
        // }
	// if percent1 >= 95 {
        //   percent1 = 1000
	// }
	fmt.Println("debug pwm1:", percent1)
      }

      if strings.Contains(input.Text(), "pwm2") {
        detectpercent2 := strings.Split(input.Text(), " ")[1]
        percent2 = getInt(detectpercent2)
        if ((percent2 == 8999) || (percent2 == 9000)) {
          percent2 = 0
        }
	// if percent2 == 0 {
        //   percent2 = 0
        // }
	// if percent2 >= 95 {
        //   percent2 = 1000
				// }
        fmt.Println("debug pwm2:", percent2)
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
    pwm(gpio.GPIO16, gpio.GPIO20)
}
