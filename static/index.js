import './form-list.js'
import './status-badge.js'

const { protocol, host } = window.location;
const socket = new WebSocket(
  `ws${protocol.includes("s") ? "s" : ""}://${host}/ws/todos`,
);

socket.onmessage = event => {
    console.log(JSON.parse(event.data))
}