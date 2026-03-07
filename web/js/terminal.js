const term = new Terminal();
term.open(document.getElementById("terminal"));

const tabs = document.getElementById("tabs");
const shellTypeSelect = document.getElementById("shellType");

let currentWS = null;

function createSession(shellType) {
    const ws = new WebSocket(`ws://${location.host}/ws`);
    ws.onopen = () => {
        ws.send(shellType);
    };

    ws.onmessage = (event) => {
        term.write(event.data);
    };

    ws.onclose = () => {
        term.write("\r\n*** Disconnected ***\r\n");
    };

    const tab = document.createElement("div");
    tab.className = "tab";
    tab.textContent = shellType;
    tab.onclick = () => {
        currentWS = ws;
        term.clear();
    };
    tabs.appendChild(tab);

    currentWS = ws;
}

document.getElementById("newSession").onclick = () => {
    createSession(shellTypeSelect.value);
};