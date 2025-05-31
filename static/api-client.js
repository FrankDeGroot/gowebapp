import * as wsc from './ws-client.js'

export default class ApiClient {
	#path

	constructor(prefix) {
		this.#path = `/api/${prefix}`
	}

	async get() {
		const response = await fetch(this.#path)
		if (!response.ok) throw "Error get"
		return await response.json()
	}

	async post(task) {
		if (wsc.wsClient.connected) {
			task.verb = "Post"
			wsc.wsClient.send(JSON.stringify(task))
		} else {
			const response = await fetch(this.#path, {
				method: "POST",
				body: JSON.stringify(task)
			})
			if (response.ok) {
				return await response.json()
			} else {
				throw "Error posting"
			}
		}
	}

	async put(task) {
		if (wsc.wsClient.connected) {
			task.verb = "Put"
			wsc.wsClient.send(JSON.stringify(task))
		} else {
			const response = await fetch(this.#path, {
				method: "PUT",
				body: JSON.stringify(task)
			})
			if (response.ok) {
				return
			} else {
				throw "Error putting"
			}
		}
	}

	async del(id) {
		const response = await fetch(`${this.#path}/${id}`, {
			method: 'DELETE'
		})
		if (!response.ok) throw "Error deleting"
	}
}
