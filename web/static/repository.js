window.addEventListener('DOMContentLoaded', event => {
	document.querySelector('#refs-menu-button').addEventListener('click', event => {
		document.querySelector('#refs-menu').classList.toggle('absolute')
		document.querySelector('#refs-menu').classList.toggle('hidden')
		event.stopPropagation()
	})

	document.querySelector('#refs-menu').addEventListener('click', event => {
		event.stopPropagation()
	})

	document.addEventListener('click', event => {
		document.querySelector('#refs-menu').classList.remove('absolute', 'hidden')
		document.querySelector('#refs-menu').classList.add('hidden')
	})
})
