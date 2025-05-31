import ApiClient from "./api-client.js"

customElements.define('form-list',
	class extends HTMLElement {
		#template = this.querySelector("template")
		#newForm = this.querySelector(".new")
		#marked = new Set()
		#prefix = this.dataset.prefix
		#apiClient = new ApiClient(this.#prefix)

		connectedCallback() {
			document.addEventListener(`${this.#prefix}:post`, e => {
				this.#toForm(e.detail)
			})
			document.addEventListener(`${this.#prefix}:put`, e => {
				this.#toForm(e.detail)
			})
			document.addEventListener(`${this.#prefix}:delete`, e => {
				const form = this.querySelector(`form[id='${this.#formId(e.detail)}'`)
				if (form) this.removeChild(form)
			})
			document.addEventListener(`${this.#prefix}:getAll`, e => {
				this.#fill(e)
			})
			this.addEventListener('submit', async e => {
				e.preventDefault()
				if (e.target === this.#newForm) {
					let response = await this.#post(this.#toObj(e.target))
					if (response) this.#toForm(response)
					e.target.reset()
				} else {
					this.#markPut(e.target)
				}
			})
			this.addEventListener('input', e => {
				const form = e.target.closest('form')
				if (form !== this.#newForm) {
					switch (e.target.type) {
						case 'text':
							this.#markPut(form)
							break
						default:
							this.#put(this.#toObj(form))
							break
					}
				}
			})
			this.addEventListener('click', async e => {
				const elm = e.target
				switch (elm.name) {
					case 'delete':
						await this.#del(form.id.value)
						try {
							this.removeChild(form)
						} catch (err) {
							if (err.name !== 'NotFoundError') throw err
						}
						break
					case 'load':
						this.#load()
						break
					default:
						const form = elm.closest('form')
						if (form === this.#newForm) return
					}
				})
			this.#load()
		}

		async #load() {
			this.#busy()
			this.#apiClient.get().then(tasks => {
				this.#done()
				if (tasks) this.#fill(tasks)
			}).catch(err => {
				this.#error()
				console.log(err)
			})
		}

		async #fill(tasks) {
			const ids = new Set()
			for (const obj of tasks) {
				ids.add(obj.id)
				this.#toForm(obj)
			}
			for (const form of this.querySelectorAll('form:not(.new)')) {
				if (!ids.has(form.id.value)) {
					this.removeChild(form)
				}
			}
		}

		#toForm(obj) {
			const form = this.#getForm(obj)
			for (const [key, value] of Object.entries(obj)) {
				const input = form.querySelector(`input[name='${key}']`)
				const label = form.querySelector(`label[for='${key}']`)
				if (label) {
					label.htmlFor = input.id = key + obj.id
				}
				switch (input.type) {
					case 'checkbox':
						input.checked = value;
					default:
						input.value = value
				}
			}
		}

		#getForm(obj) {
			const formId = this.#formId(obj)
			const form = document.getElementById(formId) ?? this.#template.content.cloneNode(true).children[0]
			if (!form.getAttribute('id')) this.insertBefore(form, this.#newForm)
			form.id = formId
			return form
		}

		#formId(obj) {
			return this.id + obj.id
		}

		#markPut(form) {
			this.#marked.add(form)
			setTimeout(() => {
				this.#marked.forEach(form => this.#put(this.#toObj(form)))
				this.#marked.clear()
			}, 1000)
		}

		async #post(obj) {
			this.#busy()
			try {
				let response = await this.#apiClient.post(obj)
				this.#done()
				return response
			} catch (err) {
				console.error(err)
				this.#error()
			}
		}

		async #put(obj) {
			this.#busy()
			try {
				await this.#apiClient.put(obj)
				this.#done()
			} catch (err) {
				console.error(err)
				this.#error()
			}
		}

		async #del(id) {
			this.#busy()
			try {
				await this.#apiClient.del(id)
				this.#done()
			} catch (err) {
				console.error(err)
				this.#error()
			}
		}

		#toObj(form) {
			const obj = {}
			for (const input of form.querySelectorAll('input')) {
				var value
				switch (input.type) {
					case 'checkbox':
						value = input.checked
					case 'button':
						break
					default:
						value = input.value
				}
				obj[input.name] = value
			}
			return obj
		}

		#busy() {
			this.#dispatchEvent('busy')
		}

		#done() {
			this.#dispatchEvent('done')
		}

		#error() {
			this.#dispatchEvent('error')
		}

		#dispatchEvent(name) {
			this.dispatchEvent(new CustomEvent(name, { bubbles: true }))
		}
	})
