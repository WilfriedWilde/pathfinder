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

const fileSelect = document.getElementById("select-file-names")
const startSelect = document.getElementById("select-start-station")
const endSelect = document.getElementById("select-end-station")
const trainSelect = document.getElementById("select-train-numbers")
const algoSelect = document.getElementById("select-algorithm")

fileSelect.name = "file"
initData.names.forEach(name => {
    const option = document.createElement("option")

    option.value = name
    option.innerText = name

    fileSelect.appendChild(option)
});