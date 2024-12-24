customElements.define('status-badge',
    class extends HTMLElement {
        connectedCallback() {
            this.innerText = '✔️'
            document.addEventListener('busy', e => {
                this.innerText = '⏲️'
            })
            document.addEventListener('done', e => {
                this.innerText = '✔️'
            })
            document.addEventListener('error', e => {
                this.innerText = '❌'
            })
        }
    })