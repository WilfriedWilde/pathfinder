async function fetchInitData() {
    try {
        const response = await fetch("http://localhost:8081/api/init")
        const data = await response.json()

        return data
    } catch(err) {
        console.log("Error fetching init data:", err)
        return []
    }
}

const initData = await fetchInitData()

const params = new URLSearchParams(window.location.search);
const selectedFile = params.get("file") || "1.txt";
let selectedStart = params.get("start");
let selectedEnd = params.get("end")
const selectedTrains = params.get("trains") || "10";
const selectedAlgos = params.getAll("algorithms");

const fileSelect = document.getElementById("select-file-names")
const startSelect = document.getElementById("select-start-station")
const endSelect = document.getElementById("select-end-station")
const trainSelect = document.getElementById("select-train-numbers")

trainSelect.value = selectedTrains;

const algoSelect = document.getElementById("select-algorithm")

console.log(initData)
initData.names.forEach(name => {
    const option = document.createElement("option")

    if (name === selectedFile) {
        option.selected = true

    }

    option.value = name
    option.innerText = name

    fileSelect.appendChild(option)
});

// init select stations
populateSelectStation(selectedFile)

initData.algorithms.forEach(algo => {
    const container = document.createElement('div')
    container.classList.add('algorithm-option-container')

    const checkbox = document.createElement('input')
    if (selectedAlgos.length > 0) {
        checkbox.checked = selectedAlgos.includes(algo);
    } else if (algo === "dijkstra") {
        checkbox.checked = true
    }

    checkbox.type = "checkbox"
    checkbox.value = algo
    checkbox.name = "algorithms"

    const label = document.createElement('label')
    label.innerText = algo

    container.append(checkbox, label)
    algoSelect.appendChild(container)
})


fileSelect.addEventListener("change", updateStations)

function updateStations(e) {
    const file = e.currentTarget.value
    
    populateSelectStation(file)
}

function populateSelectStation(file) {
    startSelect.innerHTML = ""
    endSelect.innerHTML = ""

    initData.stations[file].forEach(station => {
        const optionStart= document.createElement("option")

        optionStart.value = station
        optionStart.innerText = station

        const optionEnd = optionStart.cloneNode(true)

        startSelect.appendChild(optionStart)
        endSelect.appendChild(optionEnd)
    })

    if (selectedStart) {
        startSelect.value = selectedStart;
        selectedStart = null;
    } else {
        Array.from(startSelect.querySelectorAll('option'))[0].selected = true
    }

    if (selectedEnd) {
        endSelect.value = selectedEnd;
        selectedEnd = null;
    } else {
        const options = Array.from(endSelect.querySelectorAll('option'))
        options[options.length-1].selected = true
    }
}