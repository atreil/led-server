function sendCommand(endpoint) {
    var commandStr = `${endpoint}_command`
    var resultStr = `${endpoint}_status`

    var selector = document.getElementById(`${endpoint}_command`)
    var resultBox = document.getElementById(`${endpoint}_status`)
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

    xhttp.open("POST", `/${endpoint}`)
    xhttp.setRequestHeader("Content-Type", "application/json")
    xhttp.send(`{"Command": "${command}"}`)
}

function stopAndClear() {
    var resultBox = document.getElementById("stop_and_clear")
    var statusPrefix = "Command Status: "
    var stop = new XMLHttpRequest()
    stop.onreadystatechange = function() {
        resultBox.innerHTML = `${statusPrefix}${this.status} (${this.statusText}) ${this.responseText}`
    };

    stop.open("POST", `/daemon`)
    stop.setRequestHeader("Content-Type", "application/json")
    stop.send(`{"Command": "stop"}`)

    var clear = new XMLHttpRequest()
    clear.onreadystatechange = function() {
        resultBox.innerHTML = `${statusPrefix}${this.status} (${this.statusText}) ${this.responseText}`
    };

    clear.open("POST", `/device`)
    clear.setRequestHeader("Content-Type", "application/json")
    clear.send(`{"Command": "clear"}`)
}