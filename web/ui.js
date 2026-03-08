/* ------------------------------
   Login
--------------------------------*/
async function login() {
    const pass = document.getElementById("password").value
    const r = await fetch("/api/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
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
    const ip = document.getElementById("builderIp").value
    const port = document.getElementById("builderPort").value
    const script = document.getElementById("builderScript").value

    fetch(`/api/builder?lang=${script}&ip=${ip}&port=${port}`)
    .then(r => r.text())
    .then(cmd => {
        document.getElementById("shellOutput").textContent = cmd
    })
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
    if(terminals[id]) {
        terminals[id].tab.style.display = "block"
        terminals[id].termDiv.style.display = "block"
        activateTab(id)

        return
    }

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

    const proto = location.protocol === "https:" ? "wss" : "ws";
    const ws = new WebSocket(`${proto}://${location.host}/ws/session?id=${id}`);

    ws.onopen = ()=>term.focus()
    ws.onmessage = e=>term.write(e.data)
    term.onData(data=>{ if(ws.readyState===1) ws.send(data) })

    window.addEventListener("resize", () => {
        for (let k in terminals) {
            terminals[k].fit.fit()
        }
    })

    terminals[id]={ tab, term, termDiv, ws, fit: fitAddon }

    tab.onclick=()=>activateTab(id)
    close.onclick=(e)=>{ e.stopPropagation(); closeTab(id) }

    activateTab(id)

    term.write('=== Shell Session ===\n\rWelcome!\n\rPress [ENTER] to use openshell.\n\r');
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
    
    t.tab.style.display = "none"
    t.termDiv.style.display = "none"
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
document.addEventListener("mousemove", e => {
    if (!dragging) return

    let y = e.clientY - container.getBoundingClientRect().top
    if (y < 80) y = 80
    if (y > container.clientHeight - 120) y = container.clientHeight - 120

    sessions.style.height = y + "px"

    for (let k in terminals) {
        terminals[k].fit.fit()
    }
})

loadSessions()
setInterval(() => {
    loadSessions();
}, 3000);