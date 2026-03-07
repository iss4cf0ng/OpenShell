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

    const label=document.createElement("span")
    label.innerText=id

    const close=document.createElement("span")
    close.innerText="✕"
    close.className="tab-close"

    tab.appendChild(label)
    tab.appendChild(close)

    tabs.appendChild(tab)

    const termDiv=document.createElement("div")
    termDiv.className="terminal"
    termDiv.style.display="none"

    terminalArea.appendChild(termDiv)

    const term=new Terminal({
        cursorBlink:true
    })

    term.open(termDiv)

    const proto = location.protocol === "https:" ? "wss://" : "ws://"

    const ws = new WebSocket(proto + location.host + "/ws/session?id=" + id)

    ws.onmessage = e=>{
        term.write(e.data)
    }

    ws.onopen = ()=>{
        term.focus()
    }

    term.onData(data=>{
        if(ws.readyState === 1){
            ws.send(data)
        }
    })

    terminals[id]={
        tab,
        term,
        termDiv,
        ws
    }

    tab.onclick=()=>activateTab(id)

    close.onclick=(e)=>{
        e.stopPropagation()
        closeTab(id)
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

    if(!t) return

    t.tab.classList.add("active")
    t.termDiv.style.display="block"
}

function closeTab(id){

    const t=terminals[id]

    if(!t) return

    t.ws.close()

    t.tab.remove()
    t.termDiv.remove()

    delete terminals[id]
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