
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
    
    setInterval(updateCommandlineOutput, 400);

}

function updateCommandlineOutput() {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/get/commandline_output");
    request.onload = () => {
        // status not modified
        if (request.status == 304) {
            console.log("not modified");
            return;
        }
        console.log(request.response);
        document.getElementById("commandline-output").innerHTML = request.response;
    }
    request.send();
}

window.addEventListener("load", start, false);