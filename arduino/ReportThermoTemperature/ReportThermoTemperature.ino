#include <SoftTimer.h> // Remove Rotary and Debouncer classes, as they have another dependency
#include <SPI.h>
#include <WiFi101.h>
#include <Adafruit_MAX31856.h>
#include "secrets.h" // char ssid[], char pass[]

IPAddress server(192,168,86,142);
//IPAddress server(192,168,86,250);
//char server[] = "192.168.86.250";
int port = 1111;
int resetPort = 1112;

int resetPin = 2;

float reportTemperatureThreshold = 40.0; // Degrees C before active monitoring kicks in
unsigned long reportingDelay = 10 * 60 * 1000L; // Milliseconds between reporting
unsigned long statsDelay = 60 * 1000L; // Milliseconds between stats
unsigned long resetDelay = 60 * 60 * 1000L; // Milliseconds between hardware resets (if port is open)


// Init thermocouple with custom pins
Adafruit_MAX31856 max = Adafruit_MAX31856(5, 6, 9, 10);

void reportTemperature(Task* me);
void reportStats(Task* me);
void reset(Task* me);

// Initialize the Ethernet client library
WiFiClient client;


void setup() {
  digitalWrite(resetPin, HIGH);
  pinMode(resetPin, OUTPUT);

  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, LOW);
    //Configure pins for Adafruit ATWINC1500 Breakout
  WiFi.setPins(8,7,4);
  //Initialize serial and wait for port to open:
  Serial.begin(9600);
  while (!Serial) {
    ; // wait for serial port to connect. Needed for native USB port only
  }
  // Start thermocouple
  max.begin();
  max.setThermocoupleType(MAX31856_TCTYPE_K);

  // check for the presence of the shield:
  if (WiFi.status() == WL_NO_SHIELD) {
    Serial.println("WiFi shield not present");
    // don't continue:
    while (true);
  }

  wifiConnect();
  digitalWrite(LED_BUILTIN, HIGH);

  Task tempTask(reportingDelay, reportTemperature);
  Task statsTask(statsDelay, reportStats);
  Task resetTask(resetDelay, reset);

  SoftTimer.add(&tempTask);
  SoftTimer.add(&statsTask);
  SoftTimer.add(&resetTask);
}


void reportTemperature(Task* me) {
  float temperature;
  Serial.print(millis());

  temperature = max.readThermocoupleTemperature();
  Serial.print("Thermocouple Temp: "); Serial.println(temperature);
  if (temperature > reportTemperatureThreshold) {
    Serial.println("\nTemperature above threshold, starting connection to server... ");

    // Reconnect wifi if necessary
    if(WiFi.status() != WL_CONNECTED) {
        wifiConnect();
    }

    client.stop();
    if (client.connect(server, port)) {
      Serial.println("connected to server");


      char postData[32];
      char innerString[7];
      char outerString[7];
      dtostrf(temperature, 4, 2, innerString);
      dtostrf(max.readCJTemperature(), 4, 2, outerString);
      sprintf(postData, "inner=%s&outer=%s", innerString, outerString);
      
      client.println("POST /temperature HTTP/1.1");
      client.println("Connection: close");
      client.println("Content-Type: application/x-www-form-urlencoded");
      client.print("Content-Length: ");
      client.println(strlen(postData));
      client.println("Host: pineapplescope");
      client.println();
      client.print(postData);
      client.println();
      client.stop();
      Serial.println("Sent data");

    }
  }
}

void reportStats(Task* me) {
  // Reconnect wifi if necessary
  if(WiFi.status() != WL_CONNECTED) {
      wifiConnect();
  }

  client.stop();
  if (client.connect(server, port)) {
    Serial.println("connected to server");

    char postData[64];
    sprintf(postData, "uptime=%lu&freeMemory=%d&wifiSignal=%d", millis(), freeMemory(), WiFi.RSSI());
    Serial.println(postData);
    client.println("POST /stats HTTP/1.1");
    client.println("Connection: close");
    client.println("Content-Type: application/x-www-form-urlencoded");
    client.print("Content-Length: ");
    client.println(strlen(postData));
    client.println("Host: pineapplescope");
    client.println();
    client.print(postData);
    client.println();
    client.stop();
    Serial.println("Sent data");

  }
}

void reset(Task* me) {
    // Reconnect wifi if necessary
  if(WiFi.status() != WL_CONNECTED) {
      wifiConnect();
  }

  client.stop();
  // If we can connect on the specified port, reset the hardware
  if (client.connect(server, resetPort)) {
    digitalWrite(resetPin, LOW);
  }
}

void wifiConnect() {
  WiFi.end();
  nm_bsp_reset();

  int status = WL_IDLE_STATUS;
  // attempt to connect to WiFi network:
  while (status != WL_CONNECTED) {
    Serial.print("Attempting to connect to SSID: ");
    Serial.println(ssid);
    // Connect to WPA/WPA2 network. Change this line if using open or WEP network:
    status = WiFi.begin(ssid, pass);

    // wait 10 seconds for connection:
    delay(10000);
  }
  Serial.println("Connected to wifi");
  printWiFiStatus();
}

void printWiFiStatus() {
  // print the SSID of the network you're attached to:
  Serial.print("SSID: ");
  Serial.println(WiFi.SSID());

  // print your WiFi shield's IP address:
  IPAddress ip = WiFi.localIP();
  Serial.print("IP Address: ");
  Serial.println(ip);

  // print the received signal strength:
  long rssi = WiFi.RSSI();
  Serial.print("signal strength (RSSI):");
  Serial.print(rssi);
  Serial.println(" dBm");
}




#ifdef __arm__
// should use uinstd.h to define sbrk but Due causes a conflict
extern "C" char* sbrk(int incr);
#else  // __ARM__
extern char *__brkval;
#endif  // __arm__

int freeMemory() {
  char top;
#ifdef __arm__
  return &top - reinterpret_cast<char*>(sbrk(0));
#elif defined(CORE_TEENSY) || (ARDUINO > 103 && ARDUINO != 151)
  return &top - __brkval;
#else  // __arm__
  return __brkval ? &top - __brkval : &top - __malloc_heap_start;
#endif  // __arm__
}



