customElements.define('connected-badge',
	class extends HTMLElement {
		connectedCallback() {
			document.addEventListener('connected', _ => {
				this.innerText = '🔗'
			})
			document.addEventListener('disconnected', _ => {
				this.innerText = '⛓️‍💥'
			})
		}
	})
