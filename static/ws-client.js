export const wsClient = new class {
	#timeout = 100
	#connected = false
	#socket

	get connected() {
		return this.#connected
	}

	send(json) {
		if (this.#connected) {
			this.#socket.send(json)
		}
	}

	connect() {
		const { protocol, host } = window.location;
		this.#socket = new WebSocket(
			`ws${protocol.includes("s") ? "s" : ""}://${host}/ws`,
		);

		this.#socket.onopen = _ => {
			this.#connected = true
			document.dispatchEvent(new CustomEvent("connected", { bubbles: true }))
		}

		this.#socket.onmessage = event => {
			const e = JSON.parse(event.data)
			if (Array.isArray(e)) {
				document.dispatchEvent(new CustomEvent("tasks:getAll", { bubbles: true, detail: e }))
			} else {
				let verb;
				switch (e.verb) {
					case "Post":
						verb = "tasks:post"
						break
					case "Put":
						verb = "tasks:put"
						break
					case "Delete":
						verb = "tasks:delete"
						break
				}
				delete e.verb
				document.dispatchEvent(new CustomEvent(verb, { bubbles: true, detail: e }))
			}
		}

		this.#socket.onclose = _ => {
			this.#connected = false
			document.dispatchEvent(new CustomEvent("disconnected", { bubbles: true }))
			const wait = this.#timeout + (Math.random() - .5) * this.#timeout / 10
			console.log("Disconnected, waiting", wait, "ms to reconnect.")
			setTimeout(() => {
				this.connect();
				if (this.#timeout < 32 * this.#timeout) this.#timeout *= 2
			}, wait)
		}
	}
}

