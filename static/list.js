import { markPut } from "./put.js"
import { busy, done, error } from "./status.js"

export async function list(path, list) {
    busy()
    const response = await fetch(`/api/${path}`)
    if (response.ok) {
        done()
    } else {
        error()
    }
    setList(await response.json(), list)
}

export function setList(things, list) {
    for(const thing of things) {
        toForm(addForm(list.id + thing.id, list, list.querySelector("template")), thing)
    }
    list.addEventListener('submit', e => {
        e.preventDefault()
        markPut(e.target)
    })
    list.addEventListener('input', e => {
        markPut(e.target.closest('form'))
    })
}

function addForm(formId, list, template) {
    const form = document.getElementById(formId) ?? template.content.cloneNode(true).children[0]
    if (!form.getAttribute('id')) list.appendChild(form)
    form.id = formId
    return form
}

function toForm(form, obj) {
    for(const [key, value] of Object.entries(obj)) {
        const input = form.querySelector(`input[name='${key}']`)
        const id = input.id = key + obj.id
        const label = form.querySelector(`label[for='${key}']`)
        if(label) label.htmlFor = id
        switch(input.type) {
            case 'checkbox':
                input.checked = value;
            default:
                input.value = value
        }
    }
}