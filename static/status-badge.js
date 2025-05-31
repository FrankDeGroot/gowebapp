customElements.define('status-badge',
	class extends HTMLElement {
		connectedCallback() {
			this.innerText = '✔️'
			document.addEventListener('busy', _ => {
				this.innerText = '⏲️'
			})
			document.addEventListener('done', _ => {
				this.innerText = '✔️'
			})
			document.addEventListener('error', _ => {
				this.innerText = '❌'
			})
		}
	})
