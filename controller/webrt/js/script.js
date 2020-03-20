function sendDaemonCommand() {
    var selector = document.getElementById("daemon_command")
    var resultBox = document.getElementById("daemon_status")
    var statusPrefix = "Command Status: "
    var command = selector.options[selector.options.selectedIndex].value
    if (command === "") {
        console.debug("no option selected")
        resultBox.innerHTML = `${statusPrefix}Please select a command.`
        return
    }

    console.debug(`sending command: ${command}`)

    // just use AJAX because I'm too lazy to use a proper framework :)
    var xhttp = new XMLHttpRequest()
    xhttp.onreadystatechange = function() {
        resultBox.innerHTML = `${statusPrefix}${this.status} (${this.statusText}) ${this.responseText}`
    };

    if (command == "clear") {
        xhttp.open("POST", "/device")
    } else {
        xhttp.open("POST", "/daemon")
    }
    xhttp.setRequestHeader("Content-Type", "application/json")
    xhttp.send(`{"Command": "${command}"}`)
}