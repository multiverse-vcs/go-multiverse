<script>
	import Icon from "./Icon.svelte"
	import { path } from "../services/store"

	export let entries = []

</script>

<ul>
	{#each entries as e, i}
	<li class="px-4" >
		<a href="#" on:click={() => $path = `${$path}/${e.Name}`} class:bg-white={i % 2} class="flex items-center p-2 rounded-md">
			{#if e.IsDir}
			<span class="text-blue-500">
				<Icon name="folder" fill="currentColor" width="18" height="18" />
			</span>
			{:else}
			<span class="text-gray-400">
				<Icon name="file" width="18" height="18" />
			</span>
			{/if}
			<span class="ml-3">{e.Name}</span>
		</a>
	</li>
	{/each}
</ul>
