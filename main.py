import time
import RPi.GPIO as GPIO
import threading
import time

GPIO.setmode(GPIO.BCM)
#define mode gpio in / out
GPIO.setup(18, GPIO.OUT)
GPIO.setup(19, GPIO.OUT)
GPIO.setup(6, GPIO.IN)

class dimming1(threading.Thread):

    def __init__(self, dimming):
        self.offtime = 10000
        self.pwm = 18
        self.zc = 6
        self.dimming = 0
        self.targetdimming = dimming
        self.status = 0
        super(dimming1, self).__init__()
        GPIO.add_event_detect(self.zc, GPIO.RISING, callback=zero_crossing)

    def zero_crossing(self):
        if self.dimming == 100:
            #print("== 100")
            GPIO.output(self.pwm, GPIO.HIGH)
        elif self.dimming == 0:
            GPIO.output(self.pwm, GPIO.LOW)
        else:
            #print("dim 1 else {}".format(self.offtime/100000))
            GPIO.output(self.pwm, GPIO.HIGH)
            time.sleep(self.offtime/100000)
            GPIO.output(self.pwm, GPIO.LOW)
            self.status = self.offtime/100000


    def run(self):
        while 1:
            if self.targetdimming < self.dimming:
                self.dimming = self.dimming - 1
                self.offtime = 10000 - (100*self.dimming)
            elif self.targetdimming > self.dimming:
                self.dimming = self.dimming+1
                self.offtime = 10000 - (100*self.dimming)
            time.sleep(0.01)
            self.zero_crossing()

class dimming2(threading.Thread):

    def __init__(self, dimming):
        self.offtime = 10000
        self.pwm = 19
        self.zc = 6
        self.dimming = 0
        self.targetdimming = dimming
        self.status = 0
        super(dimming2, self).__init__()
        GPIO.add_event_detect(self.zc, GPIO.RISING, callback=zero_crossing)

    def zero_crossing(self):
        if self.dimming == 100:
            GPIO.output(self.pwm, GPIO.HIGH)
        elif self.dimming == 0:
            GPIO.output(self.pwm, GPIO.LOW)
        else:
            #print("dim 2 else {}".format(self.offtime/100000))
            GPIO.output(self.pwm, GPIO.HIGH)
            time.sleep(self.offtime/100000)
            GPIO.output(self.pwm, GPIO.LOW)
            self.status = self.offtime/100000


    def run(self):
        while 1:
            if self.targetdimming < self.dimming:
                self.dimming = self.dimming - 1
                self.offtime = 10000 - (100*self.dimming)
            elif self.targetdimming > self.dimming:
                self.dimming = self.dimming+1
                self.offtime = 10000 - (100*self.dimming)
            time.sleep(0.01)
            self.zero_crossing()

dim1, dim2 = dimming1(0).start(), dimming2(100).start()
