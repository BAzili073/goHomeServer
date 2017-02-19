package main
//GOOS=linux GOARCH=arm GOARM=6 go build
//scp ~/Work/goHomeServer/* pi@192.168.1.50:GoServer
//scp d:/Works/goserver/goserver pi@192.168.1.50:GoServer
//sudo /etc/init.d/goServer start


import (
    "fmt"
    "net/http"
    "html/template"
    "log"
    "encoding/json"
    "strconv"
    "github.com/jacobsa/go-serial/serial"
    "time"
)
var RADIO_COMMAND_RESPONSE_OK byte = 1;
var RADIO_COMMAND_PINMODE byte =  2;
var RADIO_COMMAND_DIGITAL_WRITE byte = 3;
var RADIO_COMMAND_ANALOG_WRITE byte = 4;
var RADIO_COMMAND_PIN_NOTIFY_CHANGE byte = 5;
var RADIO_COMMAND_PIN_SET_SETTING byte = 6;
var RADIO_COMMAND_PIN_RESET_SETTING byte = 7;
var RADIO_COMMAND_DIGITAL_READ byte = 8;
var RADIO_COMMAND_ANALOG_READ byte = 9;
var RADIO_COMMAND_DHT_TEMP_GET byte = 10;
var RADIO_COMMAND_DHT_TEMP_RESP byte = 11;
var RADIO_COMMAND_DHT_HUMI_GET byte = 12;
var RADIO_COMMAND_DHT_HUMI_RESP byte = 13;
var RADIO_COMMAND_DHT_ADD byte = 14;
var RADIO_COMMAND_FULL_RESET byte = 15;
var RADIO_COMMAND_CHANGE_POINT_ID byte = 16;
var RADIO_COMMAND_POINT_ON byte = 17;
var RADIO_COMMAND_PING byte = 18;
var RADIO_COMMAND_PONG byte = 19;
var RADIO_COMMAND_PING_FOR_ALL byte = 20;
var RADIO_COMMAND_CHANGE_MASTER_ID byte = 21;

var RADIO_COMMAND_DIGITALREAD_RESP byte = 108;
var RADIO_COMMAND_ANALOGREAD_RESP byte = 109;

var INPUT byte = 0; var OUTPUT byte = 1;
var LOW byte = 0;   var HIGH byte = 1;

var LOBBY_ADDRESS byte = 3;
var LOBBY_LIGHT byte = 2;

var HALL_ADDRESS byte = 1;
var HALL_RGB_LIGHT_RED byte = 5;
var HALL_RGB_LIGHT_GREEN byte = 3;
var HALL_RGB_LIGHT_BLUE byte = 6;

// "github.com/galaktor/gorf24"
var RGB_light  = map[string]int{};

type page struct {
  Title string
  Msg string
  RGB_light map[string]int
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love  %s!", r.URL.Path[1:])
}

func index(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-type","text/html");
  t, err := template.ParseFiles("./index.html")
  if err !=nil {log.Panic(err)}

  t.Execute(w, &page{Title:"Just Page",Msg: "Just Message",RGB_light : RGB_light});
}

type createValueRequest struct {
  Value string `json:"value"`
  Id string `json:"id"`
  Palette string `json:"palette_s"`
  // Value_byte byte `json:"value"`
}

type createControlRequest struct {
  Address string `json:"address"`
  Command string `json:"command"`
  Pin string `json:"pin"`
}

func controlPoints(w http.ResponseWriter, r *http.Request){
  decoder := json.NewDecoder(r.Body)
  var t createControlRequest
  err := decoder.Decode(&t)

  if err != nil {
    log.Fatal(err)
  }
  addr, err := strconv.Atoi(t.Address);
  pin, err := strconv.Atoi(t.Pin);
  if (t.Command == "cmdModeInput"){
    sendCommand([]byte{0xA9,0xA9,byte (addr),0x01,RADIO_COMMAND_PINMODE,byte (pin),INPUT});
  }else if (t.Command == "cmdModeOutput"){
    sendCommand([]byte{0xA9,0xA9,byte (addr),0x01,RADIO_COMMAND_PINMODE,byte (pin),OUTPUT});
  }else if (t.Command == "cmdHigh"){
      sendCommand([]byte{0xA9,0xA9,byte (addr),0x01,RADIO_COMMAND_DIGITAL_WRITE,byte (pin),HIGH});
  }else if (t.Command == "cmdLow"){
      sendCommand([]byte{0xA9,0xA9,byte (addr),0x01,RADIO_COMMAND_DIGITAL_WRITE,byte (pin),LOW});
  }else if (t.Command == "hall_l"){
      sendCommand([]byte{0xA9,0xA9,LOBBY_ADDRESS,0x01,RADIO_COMMAND_DIGITAL_WRITE,LOBBY_LIGHT,byte (addr)});
  }




  js, err := json.Marshal(struct{Result string `json:"result"`; Address string; Pin string; Command string }{"ok", t.Address, t.Pin, t.Command})
  // js, err := json.Marshal(struct{Result string `json:"result"`; Address byte;Color_id string}{"ok", 244,"ddd"})
  if err != nil {
    log.Fatal(err)

  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func setupPoints(w http.ResponseWriter, r *http.Request){
  sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_PINMODE,HALL_RGB_LIGHT_RED,OUTPUT});
  time.Sleep(200 * time.Millisecond);
  sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_PINMODE,HALL_RGB_LIGHT_GREEN,OUTPUT});
  time.Sleep(200 * time.Millisecond);
  sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_PINMODE,HALL_RGB_LIGHT_BLUE,OUTPUT});
  time.Sleep(200 * time.Millisecond);
  sendCommand([]byte{0xA9,0xA9,LOBBY_ADDRESS,0x01,RADIO_COMMAND_PINMODE,LOBBY_LIGHT,OUTPUT});
  js, err := json.Marshal(struct{Result string `json:"result"`}{"ok"})
  // js, err := json.Marshal(struct{Result string `json:"result"`; Address byte;Color_id string}{"ok", 244,"ddd"})
  if err != nil {
    log.Fatal(err)

  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func createValue(w http.ResponseWriter, r *http.Request){
  decoder := json.NewDecoder(r.Body)
  var t createValueRequest
  err := decoder.Decode(&t)

  if err != nil {
    log.Fatal(err)
  }

  // log.Printf("New value: %s", t.Value)


  b, err := strconv.Atoi(t.Value);
  // n, err := strconv.Atoi(t.id);
  // u:=byte(t.Id);
    RGB_light[t.Id] = b;


    if (t.Id == "off"){
      log.Printf ("Start timer on %v second",RGB_light[t.Id]);
      time.AfterFunc(time.Duration(RGB_light[t.Id]) * time.Second, func() {
           sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_RED,0x00});
           sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_BLUE,0x00});
           sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_GREEN,0x00});
      })
    }else if (t.Id == "Red" || t.Id == "Green" || t.Id == "Blue" ){ //|| t.Id == "speed" ){
      // var pin byte;
      // if (t.Id == "Red"){ pin = HALL_RGB_LIGHT_RED};
      // if (t.Id == "Green") {pin = HALL_RGB_LIGHT_GREEN};
      // if (t.Id == "Blue") {pin = HALL_RGB_LIGHT_BLUE};
      //  sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,pin,byte(RGB_light[t.Id])});

    }else if (t.Id == "palette" ){
          log.Printf("b:%v     R:%v    G:%v     B:%v",b,(b|0xFF),(b|0x00FF),(b|0x0000FF));
        // sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_RED,byte(RGB_light[t.Id] | 0xFF)});
        // sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_GREEN,byte(RGB_light[t.Id] | 0x00FF)});
        // sendCommand([]byte{0xA9,0xA9,HALL_ADDRESS,0x01,RADIO_COMMAND_ANALOG_WRITE,HALL_RGB_LIGHT_BLUE,byte(RGB_light[t.Id] | 0x0000FF)});
    }



  js, err := json.Marshal(struct{Result string `json:"result"`; Color_value int;Color_id string;Palette_s string}{"ok", RGB_light[t.Id],t.Id,t.Palette })

  if err != nil {
    log.Fatal(err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
  // Set up options.
  http.HandleFunc("/", index)
  http.HandleFunc("/api/v1/values", createValue)
  http.HandleFunc("/api/v1/control", controlPoints)
  http.HandleFunc("/api/v1/setup", setupPoints)
  http.ListenAndServe(":8080", nil)

}




func sendCommand(b []byte){
  options := serial.OpenOptions{
    PortName: "/dev/ttyAMA0",
    BaudRate: 115200,
    DataBits: 8,
    StopBits: 1,
    MinimumReadSize: 4,
  }

  // Open the port.
  port, err := serial.Open(options)
  if err != nil {
    log.Fatalf("serial.Open: %v", err)
    return
  }

  // Make sure to close it later.
  defer port.Close()

  // Write 4 bytes to the port.
  // b := []byte{0xA9,0x47,[]byte{t.Value}}

  n, err := port.Write(b)
  log.Printf("Send -> Address:%v Command: %v Arg1: %v Arg2: %v",b[2],b[4],b[5],b[6]);
  if err != nil {
    log.Fatalf("port.Write: %v", err)
  }

  fmt.Println("Wrote", n, "bytes.")
}
