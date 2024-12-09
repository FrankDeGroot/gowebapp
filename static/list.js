import { del, markPut, post } from "./put.js"
import { busy, done, error } from "./status.js"

export async function list(list) {
    const template = list.querySelector("template")
    const newForm = list.querySelector(".new")

    async function reload() {
        busy()
        const response = await fetch('/api/todos')
        if (response.ok) {
            done()
        } else {
            error()
        }
        const ids = new Set()
        for (const thing of await response.json()) {
            ids.add(thing.id)
            toForm(addForm(list.id + thing.id, list, template), thing)
        }
        for (const form of list.querySelectorAll('form:not(.new)')) {
            if (!ids.has(form.id.value)) {
                list.removeChild(form)
            }
        }
    }

    function addForm(formId, list, template) {
        const form = document.getElementById(formId) ?? template.content.cloneNode(true).children[0]
        if (!form.getAttribute('id')) list.insertBefore(form, newForm)
        form.id = formId
        return form
    }

    list.addEventListener('submit', async e => {
        e.preventDefault()
        if (e.target === newForm) {
            await post(e.target)
            e.target.reset()
            reload()
        } else {
            markPut(e.target)
        }
    })
    list.addEventListener('input', e => {
        const form = e.target.closest('form')
        if (form === newForm) return
        markPut(form)
    })
    list.addEventListener('click', async e => {
        const elm = e.target
        const form = elm.closest('form')
        if (form === newForm) return
        if (elm.name === 'delete') {
            await del(form.id.value)
            reload()
        }
    })
    reload()
}

function toForm(form, obj) {
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