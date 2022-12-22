
function start() {
    document.getElementById("start-btn").addEventListener("click",
        () => {
            fetch("/api/post/start_server", {method: "POST"}).then(
                response => console.log(response)
            );
        });
    
    document.getElementById("show-maps-btn").addEventListener("click", 
        () => {
            fetch("/api/get/maps_list", {method: "GET"}).then(
                response => console.log(response)
            );
        }
    )
    document.getElementById("show-maps-btn").addEventListener("click", 
        () => {
            fetch("/api/get/commandline_output", {method: "GET"}).then(
                response => console.log(response)
            );
        }
    )
    
}

window.addEventListener("load", start, false);