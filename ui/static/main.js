class CloseSidebar extends HTMLElement {
	constructor() {
		super();
		this.data = {};
	}

	connectedCallback() {
		this.innerHTML = `<button><img src="/static/close.svg"></button>`;
		this.data.sidebar = document.querySelector(".sidebar-l > .sidebar");
		this.querySelector("button").addEventListener("click", () => {
			this.data.sidebar.classList.add("_hide");
		});
	}

	disconnectedCallback() {
		this.innerHTML = "";
	}
}

customElements.define("close-sidebar", CloseSidebar);
