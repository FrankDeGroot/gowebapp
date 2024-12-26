customElements.define('obj-list',
    class extends HTMLElement {
        #template = this.querySelector("template")
        #newForm = this.querySelector(".new")
        #marked = new Set()
        #path = this.dataset.path

        connectedCallback() {
            this.addEventListener('submit', async e => {
                e.preventDefault()
                if (e.target === this.#newForm) {
                    this.#toForm(await this.#post(e.target))
                    e.target.reset()
                } else {
                    this.#markPut(e.target)
                }
            })
            this.addEventListener('input', e => {
                const form = e.target.closest('form')
                if (form !== this.#newForm) this.#markPut(form)
            })
            this.addEventListener('click', async e => {
                const elm = e.target
                const form = elm.closest('form')
                if (form === this.#newForm) return
                if (elm.name === 'delete') {
                    await this.#del(form.id.value)
                    this.removeChild(form)
                }
            })
            this.#reload()
        }

        async #reload() {
            this.#busy()
            const response = await fetch(this.#path)
            if (response.ok) {
                this.#done()
            } else {
                this.#error()
            }
            const ids = new Set()
            for (const obj of await response.json()) {
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
            const formId = this.id + obj.id
            const form = document.getElementById(formId) ?? this.#template.content.cloneNode(true).children[0]
            if (!form.getAttribute('id')) this.insertBefore(form, this.#newForm)
            form.id = formId
            return form
        }

        #markPut(form) {
            this.#marked.add(form)
            setTimeout(() => {
                this.#marked.forEach(form => this.#put(form))
                this.#marked.clear()
            }, 1000)
        }

        async #post(form) {
            return await this.#save('POST', form)
        }

        async #del(id) {
            this.#busy()
            const response = await fetch(`${this.#path}/${id}`, {
                method: 'DELETE'
            })
            if (response.ok) {
                this.#done()
            } else {
                this.#error()
            }
        }

        async #put(form) {
            await this.#save('PUT', form)
        }

        async #save(method, form) {
            this.#busy()
            const response = await fetch(this.#path, {
                method,
                body: JSON.stringify(this.#toObj(form))
            })
            if (response.ok) {
                this.#done()
                return await response.json()
            } else {
                this.#error()
                throw "Error saving"
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