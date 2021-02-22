<script>
  import Icon from './Icon.svelte'

  let dropdown = false
  
  function toggleDropdown() { 
    dropdown = !dropdown 
  }
</script>

<div class="relative z-10 flex-shrink-0 flex h-12" style="-webkit-app-region: drag">
  <div class="flex-1 px-4 flex justify-between">
    <div class="flex-1 flex">
 
    </div>
    <div class="ml-4 flex items-center md:ml-6">
      <!-- Profile dropdown -->
      <div class="ml-3 relative">
        <div>
          <button on:click={toggleDropdown} class="max-w-xs bg-white flex items-center text-sm text-gray-400 rounded-full focus:outline-none" id="user-menu" aria-haspopup="true">
            <span class="sr-only">Open user menu</span>
            <Icon name="settings" width="18" height="18" />
          </button>
        </div>
        {#if dropdown}
          <div class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5" role="menu" aria-orientation="vertical" aria-labelledby="user-menu">
            <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" role="menuitem">Profile</a>
            <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" role="menuitem">Settings</a>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
