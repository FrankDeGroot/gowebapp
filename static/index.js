import './form-list.js'
import './status-badge.js'

const { protocol, host } = window.location;
let timeout = 100;

function connect() {
	const socket = new WebSocket(
		`ws${protocol.includes("s") ? "s" : ""}://${host}/ws`,
	);

	console.log("Connected")

	socket.onmessage = event => {
		const e = JSON.parse(event.data)
		let verb;
		switch (e.verb) {
			case "Post":
				verb = "todo:post"
				break
			case "Put":
				verb = "todo:put"
				break
			case "Delete":
				verb = "todo:delete"
				break
		}
		delete e.verb
		document.dispatchEvent(new CustomEvent(verb, { bubbles: true, detail: e }))
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
