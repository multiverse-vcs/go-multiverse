<script>
	import { remote } from "../services/store"
	import { repo, branch, path } from "../services/store"
	import Icon from './Icon.svelte'

	let branches = []
	repo.subscribe((value) => {
		if (!value) return;
		branch.set(value.default_branch);
		branches = Object.keys(value.branches).sort();
	})

	let breadcrumbs = []
	path.subscribe((value) => {
		breadcrumbs = value.split('/').slice(1)
	})

	let branchMenu = false
	function toggleBranchMenu() {
		branchMenu = !branchMenu
	}

	function clickBreadcrumb(index) {
		path.set(`/${breadcrumbs.slice(0, index+1).join('/')}`)
	}
</script>


<div class="flex flex-1 items-center px-4 h-8">
	{#if $remote}
	<div class="flex flex-1 items-center space-x-1">
	 	<a href="#" on:click={() => $path = ''} class="text-lg font-semibold">
 		{$remote.split('/').pop()}
 		</a>
 		{#each breadcrumbs as c, i}
 		<span class="text-gray-400 text-lg">/</span>
 		<a href="#" on:click={() => clickBreadcrumb(i) } class="text-lg font-semibold">
 			{c}
 		</a>
 		{/each}
	</div>
	{/if}

	<div class="flex flex-1">
		<!-- empty space -->
	</div>

	{#if branches.length}
	<div class="relative inline-block text-left">
		<div>
			<button on:click={toggleBranchMenu} type="button" class="inline-flex space-x-1 items-center w-full rounded-full p-2 text-sm font-medium hover:bg-gray-50 focus:outline-none">
				<Icon name="git-branch" width="16" height="16" />
				<Icon name="chevron-down" width="16" height="16" />
			</button>
		</div>
		{#if branchMenu}
		<div class="absolute origin-top-right right-0 mt-2 rounded-md shadow-lg bg-white divide-y divide-gray-100" role="menu" aria-orientation="vertical">
			<div class="py-1">
				{#each branches as b}
				<a href="#" class="px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900" role="menuitem">{b}</a>
				{/each}
			</div>
		</div>
		{/if}
	</div>
	{/if}
</div>
