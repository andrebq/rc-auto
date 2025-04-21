import RPi.GPIO as GPIO #control motor board through GPIO pins
import time #set delay time to control moving distance

#If IN1=True and IN2=False right motor move forward, If IN1=False,IN2=True right motor move backward,in other cases right motor stop
IN1 = 16 #GPIO23 to IN1  right wheel direction
IN2 = 18 #GPIO24 to IN2  right wheel direction

IN3 = 13 #GPIO27
IN4 = 15 #GPIO22

GPIO.setmode(GPIO.BOARD)
GPIO.setup(IN1, GPIO.OUT)
GPIO.setup(IN2, GPIO.OUT)
GPIO.setup(IN3, GPIO.OUT)
GPIO.setup(IN4, GPIO.OUT)

class Motor:
    def __init__(self, p1, p2):
        self._p1 = p1
        self._p2 = p2
        self._enable = True

    def swap_pins(self):
        self._p1, self._p2 = self._p2, self._p1
    
    def enable(self, enable):
        self._enable = enable
        if not self._enable:
            self.stop()

    def forward(self):
        if not self._enable:
            return
        GPIO.output(self._p1, True)
        GPIO.output(self._p2, False)

    def backward(self):
        if not self._enable:
            return
        GPIO.output(self._p1, False)
        GPIO.output(self._p2, True)

    def stop(self):
        GPIO.output(self._p1, False)
        GPIO.output(self._p2, False)

right_motor = Motor(IN1, IN2)
right_motor.swap_pins()
left_motor = Motor(IN3, IN4)
left_motor.swap_pins()

right_motor.enable(True)
left_motor.enable(True)

ttl = 2

time.sleep(7)

left_motor.backward()
right_motor.backward()
print('f')
time.sleep(ttl)
left_motor.forward()
right_motor.forward()
print('b')
time.sleep(ttl)
left_motor.stop()
right_motor.stop()
print('s')

right_motor.enable(False)
left_motor.enable(False)
