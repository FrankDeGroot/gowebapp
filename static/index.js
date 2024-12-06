document.addEventListener('submit', e => {
    e.preventDefault()
    save(e.target)
})
document.addEventListener('click', e => {
    if (e.target.nodeName === 'INPUT') {
        console.log('click')
        markSave(e.target.closest('form'))
    }
})
document.addEventListener('keydown', e => {
    if (!e.repeat && e.target.nodeName === 'INPUT') {
        console.log('keydown')
        markSave(e.target.closest('form'))
    }
})

const marked = new Set();
function markSave(form) {
    marked.add(form)
    setTimeout(() => {
        marked.forEach(form => save(form))
        marked.clear()
    }, 1000)
}

async function save(form) {
    console.log('saving', form.querySelector('input[type="hidden"]').value)
    const todo = {}
    for (const input of form.querySelectorAll('input')) {
        const name = input.name
        var value
        switch(input.type) {
            case 'hidden':
                value = parseInt(input.value)
                break
            case 'checkbox':
                value = input.checked
                break
            default:
                value = input.value
        }
        todo[name] = value
    }
    const response = await fetch('/api/todos', {
        method: 'POST',
        body: JSON.stringify(todo)
    })
    if (!response.ok) {
        console.error('cannot save todo')
    }
}

async function init() {
    const response = await fetch('/api/todos')
    if (!response.ok) {
        console.error('cannot get todos')
    }
    const todos = await response.json()
    const todoList = document.getElementById('todos')
    todoTemplate = document.getElementById('todo')
    for(const todo of todos) {
        t = todoTemplate.content.cloneNode(true)
        for(const [key, value] of Object.entries(todo)) {
            const input = t.querySelector(`input[name='${key}']`)
            id = key + todo.id
            input.id = key + todo.id
            const label = t.querySelector(`label[for='${key}']`)
            if(label) label.htmlFor = id
            switch(input.type) {
                case 'checkbox':
                    input.checked = value;
                default:
                    input.value = value
            }
        }
        todoList.appendChild(t)
    }
}

init()