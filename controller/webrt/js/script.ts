// function buildPage(): void {
//     var ledBox = buildLEDBox()
//     document.body.appendChild(ledBox)
//     // var debugBox = buildDebugBox()
//     // document.appendChild(debugBox)
// }

// function buildLEDBox(): Node {
//     var tlc = document.createElement("div")
//     tlc.setAttribute("id", "led_tlc")
//     tlc.appendChild(buildVisualizationType())
//     // tlc.appendChild(buildAdvancedSettings())
//     // tlc.appendChild(buildSendCommand())
//     return tlc
// }

// function buildVisualizationType(): Node {
//     var tlc = document.createElement("div")
//     tlc.setAttribute("id", "visualization_type_tlc")

//     var label = document.createElement("label")
//     label.setAttribute("for", "led_visualization_type")
//     label.innerText = "Visualization Type:"
//     tlc.appendChild(label)

//     var select = document.createElement("select")
//     tlc.appendChild(select)

//     select.setAttribute("id", "led_visualization_type")
//     var defaultOpt = document.createElement("option")
//     defaultOpt.setAttribute("value", "")
//     defaultOpt.selected = true
//     select.appendChild(defaultOpt)
//     var options = ["Spectrum", "Spectrum Base", "Energy", "Scroll"]
//     for (let i = 0; i < options.length; i++) {
//         select.appendChild(createOption(options[i], i))
//     }

//     var spectrumBaseInput = buildInputLabel("Spectrum Base Color", "led_spectrum_base", "color")
//     spectrumBaseInput.hidden = true
//     tlc.appendChild(spectrumBaseInput)
//     select.onchange = function(ev: Event): any {
//         if (select.value === "1") {
//             spectrumBaseInput.hidden = false
//         } else {
//             spectrumBaseInput.hidden = true
//         }
//     }


//     return tlc
// }

// function buildInputLabel(name: string, id: string, type: string) {
//     var tlc = document.createElement("div")
    
//     var label = document.createElement("label")
//     label.setAttribute("for", id)
//     label.innerText = name
//     tlc.appendChild(label)

//     var input = document.createElement("input")
//     input.setAttribute("type", type)
//     input.setAttribute("id", id)
//     tlc.appendChild(input)

//     return tlc
// }

// function createOption(valueName: string, value: number): Node {
//     var opt = document.createElement("option")
//     opt.setAttribute("value", value.toString())
//     opt.innerText = valueName
//     return opt
// }

// // function buildDebugBox(): Node {
    
// // }

// buildPage()

function ledVisualizationType() {
    var select = document.getElementById("led_visualization_type") as HTMLSelectElement
    var spectrumBaseInput = document.getElementById("led_spectrum_base_opt")
    if (select.value === "1") {
        spectrumBaseInput.hidden = false
    } else {
        spectrumBaseInput.hidden = true
    }
}

function getSelectInputValue(id: string): string {
    var input = document.getElementById(id) as HTMLSelectElement
    return input.value
}

function getSelectInputValueAsInt(field: string): number | null {
    var fieldVal = getSelectInputValue(field)
    if (fieldVal === "" || fieldVal == null) {
        return null
    }
    return parseInt(fieldVal, 10)
}

function getSelectInputValueAsColor(field: string): number[] | null {
    var fieldVal = getSelectInputValue(field)
    if (fieldVal === "" || fieldVal == null) {
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
        "N_FFT_BINS": getSelectInputValueAsInt("led_n_fft_bins"),
        "MIN_FREQUENCY": getSelectInputValueAsInt("led_min_frequency"),
        "MAX_FREQUENCY": getSelectInputValueAsInt("led_max_frequency"),
        "SPECTRUM_BASE": getSelectInputValueAsColor("led_spectrum_base"),
        "VISUALIZATION_TYPE": getSelectInputValueAsInt("led_visualization_type"),
    }
    req.send(JSON.stringify(data))
}