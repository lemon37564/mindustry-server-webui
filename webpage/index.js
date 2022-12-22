
function start() {
    document.getElementById("click-btn").addEventListener("click",
        () => {
            console.log("You clicked me!");
        });
}

window.addEventListener("load", start, false);