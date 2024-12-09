import { busy, done, error } from "./status.js"

const marked = new Set()
export function markPut(form) {
    marked.add(form)
    setTimeout(() => {
        marked.forEach(form => put(form))
        marked.clear()
    }, 1000)
}

export async function post(form) {
    await save('POST', form)
}

export async function del(id) {
    busy()
    const response = await fetch(`/api/todos/${id}`, {
        method: 'DELETE'
    })
    if (response.ok) {
        done()
    } else {
        error()
    }
}

async function put(form) {
    await save('PUT', form)
}

async function save(method, form) {
    busy()
    const response = await fetch('/api/todos', {
        method,
        body: JSON.stringify(toObj(form))
    })
    if (response.ok) {
        done()
    } else {
        error()
    }
}

function toObj(form) {
    const obj = {}
    for (const input of form.querySelectorAll('input')) {
        const name = input.name
        var value
        switch (input.type) {
            case 'checkbox':
                value = input.checked
                break
            case 'button':
                break
            default:
                value = input.value
        }
        obj[name] = value
    }
    return obj
}
