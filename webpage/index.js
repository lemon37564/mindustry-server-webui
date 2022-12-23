let terminal_emu;
let wsClient;
const MAX_LINE = 1000;

function start() {
    document.getElementById("start-btn").addEventListener("click",
        () => sendCommand("host")
    );
    document.getElementById("show-maps-btn").addEventListener("click",
        () => sendCommand("maps all")
    );
    document.getElementById("runwave-btn").addEventListener("click",
        () => sendCommand("runwave")
    );
    document.getElementById("stop-btn").addEventListener("click",
        () => sendCommand("stop")
    );
    document.getElementById("gameover-btn").addEventListener("click",
        () => sendCommand("gameover")
    );

    document.getElementById("run10wave-btn").addEventListener("click",
        () => {
            for (let i = 0; i < 10; i++) {
                setTimeout(() => sendCommand("runwave"), 50 * i);
            }
        }
    );

    let btn = document.getElementById("pause-btn");
    btn.addEventListener("click",
        () => {
            if (btn.innerHTML.toLocaleLowerCase() == "pause the game") {
                btn.innerHTML = "Resume the game";
                sendCommand("pause on");
            } else {
                btn.innerHTML = "Pause the game";
                sendCommand("pause off");
            }
        }
    );
    document.getElementById("upload-map-btn").addEventListener("click", uploadMap);

    document.getElementById("kill-btn").addEventListener("click",
        () => {
            fetch("/api/post/force_restart_server", { method: "POST" })
                .then(() => { document.getElementById("commandline-output").innerHTML = ""; });
        }
    );

    document.getElementById("restart-btn").addEventListener("click",
        () => fetch("/api/post/pull_new_version_restart", { method: "POST" })
    );

    let commandInput = document.getElementById("custom-command");
    let sendCommandBtn = document.getElementById("send-custom-btn");
    // send command when press enter, and clear the input box
    commandInput.addEventListener("keyup",
        (e) => {
            if (e.key.toLocaleLowerCase() == "enter") {
                sendCommandBtn.click();
                commandInput.value = "";
            }
        })
    sendCommandBtn.addEventListener("click",
        () => {
            let command = commandInput.value;
            sendCommand(command);
        }
    );

    establishWebsocketConnection();
    terminal_emu = document.getElementById("commandline-output");
}

function establishWebsocketConnection() {
    // get current uri and combine to new websocket uri
    // result will be like ws://localhost:8086/ws/mindustry_server
    let loc = window.location, ws_uri;
    if (loc.protocol === "https:") {
        ws_uri = "wss:";
    } else {
        ws_uri = "ws:";
    }
    ws_uri += "//" + loc.host;
    ws_uri += loc.pathname + "ws/mindustry_server";
    wsClient = new WebSocket(ws_uri);
    wsClient.onmessage = onWsMessage;
}

// websocket receive message
function onWsMessage(event) {
    let data = event.data;
    data = data.replaceAll("\n", "<br>");
    terminal_emu.innerHTML += data;
    terminal_emu.scrollTo(0, terminal_emu.scrollHeight);

    // strip overflow 
    let data_arr = terminal_emu.innerHTML.split("<br>");
    if (data_arr.length > MAX_LINE) {
        data_arr = data_arr.slice(data_arr.length - MAX_LINE, data_arr.length);
        terminal_emu.innerHTML = data_arr.join("<br>");
    }
}

function sendCommand(cmd) {
    if (cmd) {
        if (wsClient.readyState == wsClient.CLOSED) {
            console.log("Reconnecting")
            establishWebsocketConnection();
            // wait until websocket connected
            setTimeout(() => wsClient.send(cmd), 200);
        } else {
            wsClient.send(cmd);
        }
    }
}

function uploadMap() {
    let input = document.createElement("input");
    input.type = "file";
    input.setAttribute("multiple", "");
    input.setAttribute("accept", ".msav");
    input.onchange = (_) => {

        for (let i = 0; i < input.files.length; i++) {
            let uploadFile = input.files[i];
            if (uploadFile) {
                let filename = uploadFile.name;
                filename = filename.replace(/.*[\/\\]/, '');

                let reader = new FileReader();

                reader.readAsArrayBuffer(uploadFile);
                reader.onload = function (e) {
                    fetch("/api/post/upload_new_map/" + filename, {
                        method: "POST", body: this.result, headers: new Headers({
                            "Content-Type": "application/octet-stream"
                        })
                    }).then(() => {
                        // when the final one was uploaded, reloadmaps
                        if (i == input.files.length - 1) {
                            setTimeout(() => sendCommand("reloadmaps"), 150);
                        }
                    });
                };
            }
        }
    };
    input.click();
}

window.addEventListener("load", start, false);