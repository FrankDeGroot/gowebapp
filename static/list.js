import { busy, done, error } from "./status.js"

class ThingList extends HTMLElement {
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
            if (form === this.#newForm) return
            this.#markPut(form)
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
        busy()
        const response = await fetch(this.#path)
        if (response.ok) {
            done()
        } else {
            error()
        }
        const ids = new Set()
        for (const thing of await response.json()) {
            ids.add(thing.id)
            this.#toForm(thing)
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
            const id = input.id = key + obj.id
            const label = form.querySelector(`label[for='${key}']`)
            if (label) label.htmlFor = id
            switch (input.type) {
                case 'checkbox':
                    input.checked = value;
                default:
                    input.value = value
            }
        }
    }

    #getForm(thing) {
        const formId = this.id + thing.id
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
        busy()
        const response = await fetch(`${this.#path}/${id}`, {
            method: 'DELETE'
        })
        if (response.ok) {
            done()
        } else {
            error()
        }
    }

    async #put(form) {
        await this.#save('PUT', form)
    }

    async #save(method, form) {
        busy()
        const response = await fetch(this.#path, {
            method,
            body: JSON.stringify(this.#toObj(form))
        })
        if (response.ok) {
            done()
            return await response.json()
        } else {
            error()
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
}
customElements.define('thing-list', ThingList)