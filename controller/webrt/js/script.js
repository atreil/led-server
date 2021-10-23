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

function getElementValue(id) {
    return document.getElementById(id).value
}

function getElementValueAsOptionalInt(field) {
    var fieldVal = getElementValue(field)
    if (fieldVal === "") {
        return null
    }
    return parseInt(fieldVal, 10)
}

function getElementValueAsOptionalColor(field) {
    var fieldVal = getElementValue(field)
    if (fieldVal === "") {
        return null
    }
    // colors have the form #rrggbb
    return [parseInt(fieldVal.slice(1, 3), 16), parseInt(fieldVal.slice(3, 5), 16), parseInt(fieldVal.slice(5, 7), 16)]
}

function sendLED() {
    var resultBox = document.getElementById("send_led_status")
    var statusPrefix = "Status: "
    var req = new XMLHttpRequest()
    req.onreadystatechange = function() {
        resultBox.innerHTML = `${statusPrefix}${this.status} (${this.statusText}) ${this.responseText}`
    };

    req.open("POST", "/led")
    req.setRequestHeader("Content-Type", "application/json")
    var data = {
        "N_FFT_BINS": getElementValueAsOptionalInt("led_n_fft_bins"),
        "MIN_FREQUENCY": getElementValueAsOptionalInt("led_min_frequency"),
        "MAX_FREQUENCY": getElementValueAsOptionalInt("led_max_frequency"),
        "SPECTRUM_BASE": getElementValueAsOptionalColor("led_spectrum_base"),
        "VISUALIZATION_TYPE": getElementValueAsOptionalInt("led_visualization_type"),
    }
    req.send(JSON.stringify(data))
}