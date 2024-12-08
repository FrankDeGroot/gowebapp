const statusDiv = document.getElementById('status')

export function busy() {
    statusDiv.innerText = '⏲️'
}

export function done() {
    statusDiv.innerText = '✔️'
}

export function error() {
    statusDiv.innerText = '❌'
}