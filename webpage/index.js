
function start() {
    document.getElementById("start-btn").addEventListener("click",
        () => {
            fetch("/api/post/start_server", { method: "POST" })
                .then(data => console.log(data));
        });

    document.getElementById("show-maps-btn").addEventListener("click",
        () => {
            fetch("/api/post/send_command", {
                method: "POST", body: JSON.stringify({ command: "maps all" }), headers: new Headers({
                    "Content-Type": "application/json"
                })
            })
        }
    )

    let interval = setInterval(() => {
        try {
            updateCommandlineOutput();
        } catch (err) {
            console.log(err);
            clearInterval(interval);
        }
    }, 400);

}

function updateCommandlineOutput() {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/get/commandline_output");
    request.onload = () => {
        if (request.status == 200) {
            // ok
            console.log(request.response);
            document.getElementById("commandline-output").innerHTML = request.response;
        } else if (request.status == 304) {
            // not modified
            console.log("not modified");
            return;
        } else {
            // deal with some error
        }
    }
    request.send();
}

window.addEventListener("load", start, false);