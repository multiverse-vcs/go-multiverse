<script>
  import NavBar from "./components/NavBar.svelte";
  import SideBar from "./components/SideBar.svelte";
  import Tree from "./components/Tree.svelte";
  import { repo, file } from "./services/store";
</script>

<div class="flex h-screen bg-gray-100 overflow-hidden">
  <SideBar />
  <div class="flex flex-1 flex-col">
    <NavBar />
    {#if $repo}
    <div class="overflow-auto">
      {#if Array.isArray($file)} 
      <Tree entries={$file} />
      {:else if typeof $file === 'string'}
      <div class="text-sm">
       {@html $file}
      </div>
      {/if}
    </div>
    {/if}
  </div>
</div>
