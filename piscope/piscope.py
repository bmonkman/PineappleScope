import sys
import logging
import time
import os

import requests
import Adafruit_GPIO
import Adafruit_DHT
# import Adafruit_MAX31856
from Adafruit_MAX31856 import MAX31856 as MAX31856
from threading import Thread, Event

logging.basicConfig(filename='pineapplescope.log', level=logging.DEBUG, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
_logger = logging.getLogger(__name__)

serverIP = '192.168.86.142'
# serverIP = '192.168.86.250:1111'
reportTemperatureThreshold = 40.0
reportingInterval = 10 * 60
statsInterval = 60

# AM2032 config
sensorType = Adafruit_DHT.AM2302
pin = 16

# MAX31856 config
software_spi = {"clk": 26, "cs": 6, "do": 19, "di": 13}
sensor = MAX31856(software_spi=software_spi)

startTime = time.time()

class ReportTemperature(Thread):
    def __init__(self, event):
        Thread.__init__(self)
        self.stopped = event

    def run(self):
        self.check()
        while not self.stopped.wait(reportingInterval):
            self.check()

    def check(self):
        innerTemp = sensor.read_temp_c()

        if innerTemp > reportTemperatureThreshold:
            ambientTemp = None
            # Log if this happens too many times?
            while ambientTemp is None:
                humidity, ambientTemp = Adafruit_DHT.read_retry(sensorType, pin)

            data = {'inner': "{:.2f}".format(innerTemp), 'outer': "{:.2f}".format(ambientTemp)}
            print(data)
            r = requests.post("http://{0}/temperature".format(serverIP), data)
            if r.status_code is not 200:
                print("Request failed: ")
                print(r)



class ReportStats(Thread):
    def __init__(self, event):
        Thread.__init__(self)
        self.stopped = event

    def run(self):
        self.check()
        while not self.stopped.wait(statsInterval):
            self.check()

    def check(self):
            innerTemp = sensor.read_temp_c()

            ambientTemp = humidity = None
            # Log if this happens too many times?
            while ambientTemp is None or humidity is None:
                humidity, ambientTemp = Adafruit_DHT.read_retry(sensorType, pin)

            data = {'temp': "{:.2f}".format(innerTemp), 'cpuTemp': getCPUtemperature(), 'freeMemory': getFreeRAM(), 'uptime': int(time.time()-startTime), 'ambientTemp': "{:.2f}".format(ambientTemp), 'humidity': "{:.2f}".format(humidity)}
            print(data)
            r = requests.post("http://{0}/stats".format(serverIP), data)
            if r.status_code is not 200:
                print("Request failed: ")
                print(r)



def getCPUtemperature():
    res = os.popen('vcgencmd measure_temp').readline()
    return(res.replace("temp=","").replace("'C\n",""))

def getFreeRAM():
    p = os.popen('free')
    i = 0
    while 1:
        i = i + 1
        line = p.readline()
        if i==2:
            return(line.split()[3])



stopFlag = Event()

temperatureThread = ReportTemperature(stopFlag)
temperatureThread.start()

statsThread = ReportStats(stopFlag)
statsThread.start()

try:
    # Main loop, just wait for an interrupt
    while True:
        time.sleep(10.0)

except KeyboardInterrupt:
    print("Ctrl-c pressed ...")
    stopFlag.set()
    sys.exit(1)


