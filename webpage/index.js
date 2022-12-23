let interval;

function start() {
    document.getElementById("start-btn").addEventListener("click",
        () => {
            fetch("/api/post/start_server", { method: "POST" })
                .then(sendCommand("host"));
        }
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

    document.getElementById("run10wave-btn").addEventListener("click",
        () => {
            for (let i = 0; i < 10; i++) {
                setTimeout(() => sendCommand("runwave"), 100 * i);
            }
        }
    );
    document.getElementById("gameover-btn").addEventListener("click",
        () => sendCommand("gameover")
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
            fetch("/api/post/kill_server", { method: "POST" })
                .then(() => { document.getElementById("commandline-output").innerHTML = ""; });
        }
    );

    document.getElementById("restart-btn").addEventListener("click",
        () => fetch("/api/post/pull_new_version_restart", { method: "POST" })
    );

    let commandInput = document.getElementById("custom-command");
    let sendCommandBtn = document.getElementById("send-custom-btn");
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

    interval = setInterval(() => { updateCommandlineOutput(false) }, 300);
    updateCommandlineOutput(true);
}

function updateCommandlineOutput(forceUpdate) {
    let request = new XMLHttpRequest();
    if (forceUpdate) {
        request.open("GET", "/api/get/commandline_output?force_update=true");
    } else {
        request.open("GET", "/api/get/commandline_output");
    }
    request.onload = () => {
        if (request.status == 200) {
            // ok
            let response = request.response;
            response = response.replaceAll("\n", "<br>");
            document.getElementById("commandline-output").innerHTML = response;
        } else if (request.status == 304) {
            // not modified
            return;
        } else {
            // TODO: deal with some error
        }
    }
    request.send();
}

function sendCommand(cmd) {
    if (cmd) {
        fetch("/api/post/send_command", {
            method: "POST", body: JSON.stringify({ command: cmd }), headers: new Headers({
                "Content-Type": "application/json"
            })
        });
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