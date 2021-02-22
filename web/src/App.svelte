<script>
  import NavBar from './components/NavBar.svelte'
  import Author from './pages/Author.svelte'
  import Repo from './pages/Repo.svelte'
</script>

<div class="h-screen flex overflow-hidden bg-gray-100">
  <div class="flex w-0 flex-1 overflow-hidden">
    <div class="flex-1 min-w-0 min-h-full flex">
      <div class="flex-shrink-0 w-64 border-r border-gray-200 bg-white">
        <div class="h-12 w-full" style="-webkit-app-region: drag"></div>
        <Author />
      </div>
      <div class="min-w-0 flex-1">
        <NavBar />
        <Repo />
      </div>
    </div>
  </div>
</div>