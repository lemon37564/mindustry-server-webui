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
    document.getElementById("gameover-btn").addEventListener("click",
        () => sendCommand("gameover")
    );

    let btn = document.getElementById("pause-btn");
    btn.addEventListener("click",
        () => {
            if (btn.innerHTML == "Game pause") {
                btn.innerHTML = "Game resume";
                sendCommand("pause on");
            } else {
                btn.innerHTML = "Game pause";
                sendCommand("pause off");
            }
        }
    );
    document.getElementById("send-custom-btn").addEventListener("click",
        () => {
            let command = document.getElementById("custom-command").value;
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
    fetch("/api/post/send_command", {
        method: "POST", body: JSON.stringify({ command: cmd }), headers: new Headers({
            "Content-Type": "application/json"
        })
    });
}

window.addEventListener("load", start, false);