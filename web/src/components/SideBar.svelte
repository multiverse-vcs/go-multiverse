<script>
  import Profile from './Profile.svelte'
  import Icon from './Icon.svelte'
  import { self, remote } from '../services/store'

  let repos = []
  self.subscribe(value => {
    if (!value) return
    repos = Object.keys(value.author.repositories).sort()
  })
</script>

<div class="flex-shrink-0 w-64 border-r border-gray-200 bg-white">
  <div class="h-6 w-full" style="-webkit-app-region: drag">
    <!-- draggable area -->
  </div>
  <div class="pl-4 pr-6" >
    <div class="flex items-center justify-between">
      <div class="flex-1 space-y-8">
        <div class="space-y-8 block space-y-8">
          {#if $self}
            <Profile peerID={$self.peerID} />
          {/if}
          <!-- Action buttons -->
          <div class="flex flex-col">
            <button type="button" class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none">
              <span class="ml-1">New Repository</span>
            </button>
          </div>
        </div>
        <!-- Meta info -->
        <div class="flex flex-col space-y-3">
          <div class="flex items-center space-x-2">
            <span class="text-xs text-gray-500 font-medium">Repositories</span>
          </div>
          {#each repos as name}
          <div class="flex items-center space-x-2">
            <Icon name="book" width="16" height="16" />
            <a href="#" class="w-full text-md" on:click={() => $remote = `${$self.peerID}/${name}`}>{name}</a>
          </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>