import sys
import logging
import time
import os

import requests
import Adafruit_GPIO
import Adafruit_DHT
# import Adafruit_MAX31856
# from Adafruit_MAX31856 import MAX31856 as MAX31856
from MAX31856 import max31856
from threading import Thread, Event

logging.basicConfig(filename='pineapplescope.log', level=logging.DEBUG, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
_logger = logging.getLogger(__name__)

serverIP = '192.168.86.142:1111'
# serverIP = '192.168.86.250:1111'

reportTemperatureThreshold = 40.0
reportingInterval = 10 * 60
statsInterval = 60
thermocoupleTemperatureModifier = 1.0
# 1.02 at 962 (962 TC vs 978 kiln)

# AM2302 config
sensorType = Adafruit_DHT.AM2302
pin = 14

# MAX31856 config
# software_spi = {"clk": 17, "do": 4, "di": 3, "cs": 2}
# max = MAX31856(software_spi=software_spi)
max = max31856.max31856(2,4,3,17)

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
        try:
            # innerTemp = max.read_temp_c()
            innerTemp = max.readThermocoupleTemp()*thermocoupleTemperatureModifier
        except FaultError as e:
            print("Exception:")
            print(e)
            return

        if innerTemp > reportTemperatureThreshold:
            humidity, ambientTemp = Adafruit_DHT.read_retry(sensorType, pin, 5, 2)

            if ambientTemp is None:
                ambientTemp = 0.0
                print("Couldn't read ambientTemp for ReportTemperature")

            data = {'inner': "{:.2f}".format(innerTemp), 'outer': "{:.2f}".format(ambientTemp)}
            print(data)

            try:
                r = requests.post("http://{0}/temperature".format(serverIP), data)
                if r.status_code is not 200:
                    print("Request failed: ")
                    print(r)
            except requests.exceptions.RequestException as e:
                print("Request failed: ")
                print(e)


class ReportStats(Thread):
    def __init__(self, event):
        Thread.__init__(self)
        self.stopped = event

    def run(self):
        self.check()
        while not self.stopped.wait(statsInterval):
            self.check()

    def check(self):
        try:
            # innerTemp = max.read_temp_c()
            innerTemp = max.readThermocoupleTemp()*thermocoupleTemperatureModifier
        except FaultError as e:
            print("Exception:")
            print(e)
            return

        humidity, ambientTemp = Adafruit_DHT.read_retry(sensorType, pin, 5, 2)

        if humidity is None:
            humidity = 0.0
            print("Couldn't read humidity for ReportStats")

        if ambientTemp is None:
            ambientTemp = 0.0
            print("Couldn't read ambientTemp for ReportStats")

        data = {'temp': "{:.2f}".format(innerTemp), 'cpuTemp': getCPUtemperature(), 'freeMemory': getFreeRAM(), 'uptime': int(time.time()-startTime), 'ambientTemp': "{:.2f}".format(ambientTemp), 'humidity': "{:.2f}".format(humidity)}
        print(data)

        try:
            r = requests.post("http://{0}/stats".format(serverIP), data)
            if r.status_code is not 200:
                print("Stats Request failed: ")
                print(r)
        except requests.exceptions.RequestException as e:
            print("Stats Request failed: ")
            print(e)



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


