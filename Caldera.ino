#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <OneWire.h>
#include <DallasTemperature.h>

#ifndef STASSID
#define STASSID "awesome"
#define STAPSK  ""
#endif
#define CONTROLLER_PIN 15
#define DS18S20_PIN 14

const char *ssid = STASSID;
const char *password = STAPSK;
char hostString[16] = {0};
ESP8266WebServer server(80);
OneWire oneWire(DS18S20_PIN);
DallasTemperature sensors(&oneWire);

void handleRelay() {
  String state = server.arg("state");
  bool on = state == "on";
  if(on){
    digitalWrite(CONTROLLER_PIN, HIGH);
  }else{
    digitalWrite(CONTROLLER_PIN, LOW);   
  }
  server.send(200, "text/html");
}

void handleUpTime() {
  char uptime[128];
  int sec = millis() / 1000;
  int min = sec / 60;
  int hr = min / 60;
  snprintf(uptime, 128,"%02d:%02d:%02d", hr, min % 60, sec % 60);
  server.send(200, "text/html", uptime);
}

void handleTemp() {
  char temp[128];
  sensors.requestTemperatures(); 
  float temperatureC = sensors.getTempCByIndex(0);
  snprintf(temp, 128,"%f", temperatureC);
  server.send(200, "text/html", temp);
}

void handleNotFound() {
  String message = "File Not Found\n\n";
  message += "URI: ";
  message += server.uri();
  message += "\nMethod: ";
  message += (server.method() == HTTP_GET) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server.args();
  message += "\n";

  for (uint8_t i = 0; i < server.args(); i++) {
    message += " " + server.argName(i) + ": " + server.arg(i) + "\n";
  }

  server.send(404, "text/plain", message);
}

void setup(void) {
  pinMode(LED_BUILTIN, OUTPUT);
  pinMode(CONTROLLER_PIN, OUTPUT);

  //Serial.begin(115200);
  sprintf(hostString, "comedor");
  //Serial.print("Hostname: ");
  //Serial.println(hostString);
  WiFi.hostname(hostString);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  //Serial.println("");

  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    //Serial.print(".");
  }

  //Serial.println("");
  //Serial.print("Connected to ");
  //Serial.println(ssid);
  //Serial.print("IP address: ");
  //Serial.println(WiFi.localIP());



  if (MDNS.begin(hostString)) {
    //Serial.println("MDNS responder started");
  }
  MDNS.addService("http", "tcp", 80);

  server.on("/uptime", handleUpTime);  
  server.on("/temp", handleTemp);
  server.on("/relay", handleRelay);
  server.on("/inline", []() {
    server.send(200, "text/plain", "this works as well");
  });
  server.onNotFound(handleNotFound);
  server.begin();
  //Serial.println("HTTP server started");
  digitalWrite(CONTROLLER_PIN, LOW);

  sensors.begin();
  //Serial.println("Temperature sensor started");

}

void loop(void) {
  server.handleClient();
  MDNS.update();
}
