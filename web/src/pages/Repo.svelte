<script>
	import Icon from "../components/Icon.svelte"
	import Tree from "../components/Tree.svelte"
	import { file } from "../services/store"
</script>

<div class="flex-grow w-full h-full overflow-auto mt-3">
	{#if Array.isArray($file)} 
	<Tree entries={$file} />
	{:else if typeof $file === 'string'}
  <article class="text-sm">
	 {@html $file}
  </article>
	{/if}
</div>
