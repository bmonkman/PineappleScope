
#include <SPI.h>
#include <WiFi101.h>
#include <Adafruit_MAX31856.h>
#include "secrets.h"

int keyIndex = 0; // network key Index number (needed only for WEP)

int status = WL_IDLE_STATUS;
IPAddress server(192,168,86,142);
//IPAddress server(192,168,86,250);
//char server[] = "192.168.86.250";
int port = 1111;

float reportTemperatureThreshold = 40.0; // Degrees C before active monitoring kicks in
unsigned long reportingDelay = 5 * 60 * 1000L; // Milliseconds between reporting


// Initialize the Ethernet client library
WiFiClient client;

// Init thermocouple with custom pins
Adafruit_MAX31856 max = Adafruit_MAX31856(5, 6, 9, 10);

void setup() {
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
  digitalWrite(LED_BUILTIN, HIGH);
}

float temperature;
long currentMillis = reportingDelay*-1;
void loop() {
  Serial.print(millis()); Serial.print(" > "); Serial.println(currentMillis + reportingDelay);
  if (millis() > currentMillis + reportingDelay)
  {
    currentMillis = millis();
    temperature = max.readThermocoupleTemperature();
    Serial.print("Thermocouple Temp: "); Serial.println(temperature);
    if (temperature > reportTemperatureThreshold) {
      Serial.println("\nTemperature above threshold, starting connection to server... ");
  
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
  delay(1000);
  digitalWrite(LED_BUILTIN, HIGH);
  delay(1000);
  digitalWrite(LED_BUILTIN, LOW);
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





