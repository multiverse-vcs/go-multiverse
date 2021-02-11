window.addEventListener('DOMContentLoaded', event => {
 	// Enable clipboard copy on click.
  document.querySelectorAll('[data-copy-text]').forEach(el => {
		el.addEventListener('click', async (event) => {
  		const text = el.getAttribute('data-copy-text')
    	await navigator.clipboard.writeText(text)
    })
	})

	document.querySelector('#mobile-menu-button').addEventListener('click', event => {
		document.querySelector('#mobile-menu').classList.toggle('hidden')
		document.querySelector('#mobile-menu').classList.toggle('block')
		document.querySelector('#mobile-menu-open').classList.toggle('hidden')
		document.querySelector('#mobile-menu-open').classList.toggle('block')
		document.querySelector('#mobile-menu-close').classList.toggle('hidden')
		document.querySelector('#mobile-menu-close').classList.toggle('block')
	})

	document.querySelector('#user-menu-button').addEventListener('click', event => {
		document.querySelector('#user-menu').classList.toggle('absolute')
		document.querySelector('#user-menu').classList.toggle('hidden')
		event.stopPropagation()
	})

	document.querySelector('#user-menu').addEventListener('click', event => {
		event.stopPropagation()
	})

	document.addEventListener('click', event => {
		document.querySelector('#user-menu').classList.remove('absolute', 'hidden')
		document.querySelector('#user-menu').classList.add('hidden')
	})

	document.querySelector('#search-form').addEventListener('submit', event => {
		event.target.action = '/' + document.querySelector('#search').value
	})
})
