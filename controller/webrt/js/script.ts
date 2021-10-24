function buildPage(): void {
    var ledBox = buildLEDBox()
    document.body.appendChild(ledBox)
    // var debugBox = buildDebugBox()
    // document.appendChild(debugBox)
}

function buildLEDBox(): Node {
    var tlc = document.createElement("div")
    tlc.setAttribute("id", "led_tlc")
    tlc.appendChild(buildVisualizationType())
    // tlc.appendChild(buildAdvancedSettings())
    // tlc.appendChild(buildSendCommand())
    return tlc
}

function buildVisualizationType(): Node {
    var tlc = document.createElement("div")
    tlc.setAttribute("id", "visualization_type_tlc")

    var label = document.createElement("label")
    label.setAttribute("for", "led_visualization_type")
    label.innerText = "Visualization Type:"
    tlc.appendChild(label)

    var select = document.createElement("select")
    tlc.appendChild(select)

    select.setAttribute("id", "led_visualization_type")
    var defaultOpt = document.createElement("option")
    defaultOpt.setAttribute("value", "")
    defaultOpt.selected = true
    select.appendChild(defaultOpt)
    var options = ["Spectrum", "Spectrum Base", "Energy", "Scroll"]
    for (let i = 0; i < options.length; i++) {
        select.appendChild(createOption(options[i], i))
    }

    var spectrumBaseInput = buildInputLabel("Spectrum Base Color", "led_spectrum_base", "color")
    spectrumBaseInput.hidden = true
    tlc.appendChild(spectrumBaseInput)
    select.onchange = function(ev: Event): any {
        if (select.value === "1") {
            spectrumBaseInput.hidden = false
        } else {
            spectrumBaseInput.hidden = true
        }
    }


    return tlc
}

function buildInputLabel(name: string, id: string, type: string) {
    var tlc = document.createElement("div")
    
    var label = document.createElement("label")
    label.setAttribute("for", id)
    label.innerText = name
    tlc.appendChild(label)

    var input = document.createElement("input")
    input.setAttribute("type", type)
    input.setAttribute("id", id)
    tlc.appendChild(input)

    return tlc
}

function createOption(valueName: string, value: number): Node {
    var opt = document.createElement("option")
    opt.setAttribute("value", value.toString())
    opt.innerText = valueName
    return opt
}

// function buildDebugBox(): Node {
    
// }

buildPage()