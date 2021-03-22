import time
import RPi.GPIO as GPIO
import threading
import time

GPIO.setmode(GPIO.BCM)
#define mode gpio in / out
GPIO.setup(18, GPIO.OUT)
GPIO.setup(19, GPIO.OUT)
GPIO.setup(6, GPIO.IN)

class PWM1(threading.Thread):

    def __init__(self, dimming):
        self.offtime = 10000
        self.pwm = 18
        self.zc = 6
        self.dutycycle = 0
        self.frequency = dimming
        self.status = 0
        GPIO.add_event_detect(self.zc, GPIO.RISING, callback=self.zero_crossing)
        super(PWM1, self).__init__()

    def zero_crossing(self):
        if self.dutycycle == 100:
            #print("== 100")
            GPIO.output(self.pwm, GPIO.HIGH)
        elif self.dutycycle == 0:
            GPIO.output(self.pwm, GPIO.LOW)
        else:
            #print("dim 1 else {}".format(self.offtime/100000))
            GPIO.output(self.pwm, GPIO.HIGH)
            time.sleep(self.offtime/100000)
            GPIO.output(self.pwm, GPIO.LOW)
            self.status = self.offtime/100000

    def run(self):
        while 1:
            if self.frequency < self.dimming:
                self.dutycycle = self.dutycycle - 1
                self.offtime = 10000 - (100*self.dutycycle)
            elif self.frequency > self.dutycycle:
                self.dutycycle = self.dutycycle+1
                self.offtime = 10000 - (100*self.dutycycle)
            time.sleep(0.01)
            if self.status == "stop":
                break
            #self.zero_crossing()

class PWM2(threading.Thread):

    def __init__(self, dimming):
        self.offtime = 10000
        self.pwm = 19
        self.zc = 6
        self.dutycycle = 0
        self.frequency = dimming
        self.status = 0
        GPIO.add_event_detect(self.zc, GPIO.RISING, callback=self.zero_crossing)
        super(PWM2, self).__init__()

    def zero_crossing(self):
        if self.dutycycle == 100:
            GPIO.output(self.pwm, GPIO.HIGH)
        elif self.dutycycle == 0:
            GPIO.output(self.pwm, GPIO.LOW)
        else:
            #print("dim 2 else {}".format(self.offtime/100000))
            GPIO.output(self.pwm, GPIO.HIGH)
            time.sleep(self.offtime/100000)
            GPIO.output(self.pwm, GPIO.LOW)
            self.status = self.offtime/100000

    def run(self):
        while 1:
            if self.frequency < self.dutycycle:
                self.dutycycle = self.dutycycle - 1
                self.offtime = 10000 - (100*self.dutycycle)
            elif self.frequency > self.dutycycle:
                self.dutycycle = self.dutycycle+1
                self.offtime = 10000 - (100*self.dutycycle)
            time.sleep(0.01)
            if self.status == "stop":
                break
            #self.zero_crossing()

# initial by frequency
PWM1, PWM2 = PWM1(0), PWM2(0)
PWM1.start()
time.sleep(1)
# change frequency 60 hz
PWM1.frequency = 60
# change dutycycle 10%
PWM1.dutycycle = 10
# stop process
PWM1.status = "stop"
