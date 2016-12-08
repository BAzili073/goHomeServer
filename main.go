package main

import (
    "fmt"
    "net/http"
    "html/template"
<<<<<<< HEAD
    "log"
    "encoding/json"
    "strconv"
    "github.com/jacobsa/go-serial/serial"
)

 var RGB_light  = map[string]int{};


//GOOS=linux GOARCH=arm GOARM=6 go build
  //scp ~/Work/goHomeServer/* pi@192.168.1.50:GoServer
    //scp d:/Works/goserver/goserver pi@192.168.1.50:GoServer
=======
)

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
  c:=byte(b);
  // n, err := strconv.Atoi(t.id);
  // u:=byte(t.Id);
    RGB_light[t.Id] = int(c);
    log.Printf("%c =  %v",t.Id[0], RGB_light[t.Id])
  sendCommand([]byte{0xA9,t.Id[0],c});


  js, err := json.Marshal(struct{Result string `json:"result"`; Color_value byte;Color_id string}{"ok", c,t.Id })

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
