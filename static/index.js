import './form-list.js'
import './status-badge.js'

const { protocol, host } = window.location;
// TODO recreate connection if broken (onclose?)
const socket = new WebSocket(
  `ws${protocol.includes("s") ? "s" : ""}://${host}/ws/todos`,
);

socket.onmessage = event => {
    const e = JSON.parse(event.data)
    let action;
    switch(e.action) {
      case "A":
        action = "todo:add"
        break
      case "C":
        action = "todo:change"
        break
      case "D":
        action = "todo:delete"
        break
    }
    delete e.action
    document.dispatchEvent(new CustomEvent(action, { bubbles: true, detail: e }))
}
