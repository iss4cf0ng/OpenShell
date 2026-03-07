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

    if(terminals[id]){
        activateTab(id)
        return
    }

    const tab=document.createElement("div")
    tab.className="tab"
    tab.innerText=id
    tab.onclick=()=>activateTab(id)

    tabs.appendChild(tab)

    const termDiv=document.createElement("div")
    termDiv.className="terminal"
    termDiv.style.display="none"

    terminalArea.appendChild(termDiv)

    const term=new Terminal({
        cursorBlink:true
    })

    term.open(termDiv)

    term.write("Connected to "+id+"\r\n")

    /* websocket */

    const ws = new WebSocket(`ws://${location.host}/ws/session?id=${id}`)

    ws.onmessage = function(e){
        term.write(e.data)
    }

    term.onData(function(data){
        ws.send(data)
    })

    terminals[id]={
        tab,
        term,
        termDiv,
        ws
    }

    activateTab(id)
}

function activateTab(id){

    for(let k in terminals){

        const t=terminals[k]

        t.tab.classList.remove("active")
        t.termDiv.style.display="none"
    }

    const t=terminals[id]

    t.tab.classList.add("active")
    t.termDiv.style.display="block"
}

function loadSessions(){
    fetch("/api/sessions")
    .then(r=>r.json())
    .then(list=>{
        sessionTable.innerHTML = `
        <tr>
        <th>ID</th>
        <th>IP</th>
        <th>Type</th>
        </tr>
        `

        list.forEach(s=>{
            addSession(s.id,s.ip,s.type)
        })

    })
}

/* splitter */

const splitter=document.getElementById("splitter")
const sessions=document.getElementById("sessions")
const container=document.getElementById("container")

let dragging=false

splitter.addEventListener("mousedown",()=>{
    dragging=true
})

document.addEventListener("mouseup",()=>{
    dragging=false
})

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