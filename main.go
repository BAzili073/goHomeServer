package main

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

 var RGB_light  = map[string]int{};


//GOOS=linux GOARCH=arm GOARM=6 go build
  //scp ~/Work/goHomeServer/* pi@192.168.1.50:GoServer
    //scp d:/Works/goserver/goserver pi@192.168.1.50:GoServer

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
  // Value_byte byte `json:"value"`
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
          sendCommand([]byte{0xA9,'R',0x00});
          sendCommand([]byte{0xA9,'G',0x00});
          sendCommand([]byte{0xA9,'B',0x00});
      })
    }else if (t.Id == "Red" || t.Id == "Green" || t.Id == "Blue" || t.Id == "speed" ){
       sendCommand([]byte{0xA9,t.Id[0],byte(RGB_light[t.Id])});
       log.Printf("%c =  %v",t.Id[0], RGB_light[t.Id]);
    }



  js, err := json.Marshal(struct{Result string `json:"result"`; Color_value int;Color_id string}{"ok", RGB_light[t.Id],t.Id })

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
    http.ListenAndServe(":8080", nil)
}




func sendCommand(b []byte){
  options := serial.OpenOptions{
    PortName: "/dev/ttyAMA0",
    BaudRate: 9600,
    DataBits: 8,
    StopBits: 1,
    MinimumReadSize: 4,
  }

  // Open the port.
  port, err := serial.Open(options)
  if err != nil {
    log.Fatalf("serial.Open: %v", err)
  }

  // Make sure to close it later.
  defer port.Close()

  // Write 4 bytes to the port.
  // b := []byte{0xA9,0x47,[]byte{t.Value}}
  n, err := port.Write(b)
  if err != nil {
    log.Fatalf("port.Write: %v", err)
  }

  fmt.Println("Wrote", n, "bytes.")
}
