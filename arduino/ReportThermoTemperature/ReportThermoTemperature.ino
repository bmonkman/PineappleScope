
#include <SPI.h>
#include <WiFi101.h>
#include <Adafruit_MAX31856.h>
#include "secrets.h"

int keyIndex = 0; // network key Index number (needed only for WEP)

int status = WL_IDLE_STATUS;
IPAddress server(192,168,86,142);
//char server[] = "192.168.86.250";
int port = 8177;

float reportTemperatureThreshold = 40.0; // Degrees C before active monitoring kicks in
int overTemperatureDelay = 5; // Seconds between reporting when actively monitoring
int underTemperatureDelay = 20; // Seconds between checks when not actively monitoring


// Initialize the Ethernet client library
WiFiClient client;

// Init thermocouple with custom pins
Adafruit_MAX31856 max = Adafruit_MAX31856(5, 6, 9, 10);

void setup() {
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

  Serial.println("\nStarting connection to server...");
  // if you get a connection, report back via serial:
}

float temperature;
void loop() {

  temperature = max.readThermocoupleTemperature();
  if (temperature > reportTemperatureThreshold) {

    if (client.connect(server, port)) {
      Serial.println("connected to server");

      Serial.print("Thermocouple Temp: "); Serial.println(temperature);

      client.print("POST /reportData?innerTemperature=");
      client.print(temperature);
      client.print("&outerTemperature=");
      client.println(max.readCJTemperature());    
      client.println("Connection: close");
      client.println();
      client.stop();
    }
    delay(overTemperatureDelay*1000);  
  } else {
    delay(underTemperatureDelay*1000);
  }
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





