# RobotDIN_Raspberrypi_Triac
TRIAC two channel control PWM over raspberry pi Golang

![run program](https://i.ibb.co/kgsmbGD/Capture-du-2021-03-28-15-54-16.png)
![run program](https://i.ibb.co/8jFtYVJ/Capture-du-2021-03-28-15-57-19.png)

## install
```
go get github.com/warthog618/gpio
```

## run 
```
go run main.go
```

## compilation go
```
go build
./binary
```

## Command percent
the script accept to 2 float number in percent  
```
pwm1 10.8  
pwm1 10.99  
pwm1 80.75
or 
pwm2 80
pwm2 65.6
pwm2 10.87
```

## GPIO config
the script run on zero crossing pin 6  
pin 18 [pwm1] and 19 [pwm2]

Enjoy
