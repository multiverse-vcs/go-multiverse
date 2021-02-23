<script>
  import NavBar from "./components/NavBar.svelte";
  import SideBar from "./components/SideBar.svelte";
  import Repo from "./pages/Repo.svelte";
  import { repo } from "./services/store";
</script>

<div class="h-screen flex overflow-hidden bg-gray-100">
  <div class="flex w-0 flex-1">
    <div class="flex-1 min-w-0 min-h-full flex">
      <SideBar />
      <div class="min-w-0 flex-1">
        <div class="h-6 w-full" style="-webkit-app-region: drag">
          <!-- draggable area -->
        </div>
        <NavBar />
        {#if $repo}
        <Repo />
        {/if}
      </div>
    </div>
  </div>
</div>
