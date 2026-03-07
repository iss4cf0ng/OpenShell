var term = new Terminal()

term.open(document.getElementById("terminal"))

var ws = new WebSocket("ws://localhost:8080/ws")

ws.onmessage = function(event){
    term.write(event.data)
}

term.onData(function(data){
    ws.send(data)
})