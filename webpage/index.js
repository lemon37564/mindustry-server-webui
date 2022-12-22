
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
            }).then(() => {
                getCommandlineOutput();
            })
        }
    )
    document.getElementById("get-output-btn").addEventListener("click",
        () => {
            fetch("/api/get/commandline_output", { method: "GET" }).then(
                response => console.log(response)
            );
        }
    )

}

function getCommandlineOutput() {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/get/commandline_output");
    request.onload = () => {
        document.getElementById("commandline-output").innerHTML = request.response;
    }
    request.send();
}

window.addEventListener("load", start, false);