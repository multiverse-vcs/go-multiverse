<script>
  import { onDestroy } from 'svelte'
  import Profile from '../components/Profile.svelte'
  import Icon from '../components/Icon.svelte'
  import { author, authorID, peerID, repoName } from '../services/store'
  
  let names = []

  const unsubscribe = author.subscribe(value => {
    if (!value) return

    names = Object.keys(value.repositories).sort()
  })

  onDestroy(unsubscribe)
</script>

<div class="pl-4 pr-6">
  <div class="flex items-center justify-between">
    <div class="flex-1 space-y-8">
      <div class="space-y-8 block space-y-8">
        {#if $authorID}
          <Profile peerID={$authorID} />
        {/if}
        <!-- Action buttons -->
        <div class="flex flex-col">
          {#if $peerID === $authorID}
            <button type="button" class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none">
              <span class="ml-1">New Repository</span>
            </button>
          {:else}
            <button type="button" class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none">
              <span class="ml-1">Follow</span>
            </button>
          {/if}
        </div>
      </div>
      <!-- Meta info -->
      <div class="flex flex-col space-y-3">
        <div class="flex items-center space-x-2">
          <span class="text-xs text-gray-500 font-medium">Repositories</span>
        </div>
        {#each names as name}
        <div class="flex items-center space-x-2">
          <Icon name="book" width="16" height="16" />
          <a href="#" class="w-full text-md" on:click={() => $repoName = name}>{name}</a>
        </div>
        {/each}
      </div>
    </div>
  </div>
</div>
