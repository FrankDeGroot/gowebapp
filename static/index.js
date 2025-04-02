import './form-list.js'
import './status-badge.js'

const { protocol, host } = window.location;
let timeout = 100;

function connect() {
  const socket = new WebSocket(
    `ws${protocol.includes("s") ? "s" : ""}://${host}/ws/todos`,
  );

  console.log("Connected")

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

  socket.onclose = event => {
    socket.onmessage = null
    socket.onclose = null
    const wait = timeout + (Math.random() - .5) * timeout / 10
    console.log("Disconnected, waiting", wait, "ms to reconnect.")
    setTimeout(() => {
      console.log("Reconnecting")
      connect();
      if (timeout < 32 * timeout) timeout *= 2
    }, wait)
  }
}

connect()
