window.addEventListener('DOMContentLoaded', (event) => {
    document.querySelectorAll('[data-copy-text]').forEach(el => {
    	const text = el.getAttribute('data-copy-text')
		el.addEventListener('click', async (event) => {
	    	await navigator.clipboard.writeText(text)
	    })
	})
})
