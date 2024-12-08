import { list } from './list.js'
import { post } from './put.js'

function init() {
    document.getElementById("addTodo").addEventListener("submit", async e => {
        e.preventDefault()
        await post(e.target)
        e.target.reset()
        list('todos', document.getElementById('todos'))
    })
    list('todos', document.getElementById('todos'))
}

init()