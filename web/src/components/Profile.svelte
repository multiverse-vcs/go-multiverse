<script>
  export let peerID

  function shortPeerID(id) {
    if (!id) return ''
      
    return `${id.slice(0, 4)}..${id.slice(-5, -1)}`
  }
</script>

<div class="flex items-center space-x-3">
  <span class="flex-shrink-0 h-12 w-12 rounded-full overflow-hidden bg-gray-100">
    <svg class="h-full w-full text-gray-300" fill="currentColor" viewBox="0 0 24 24">
      <path d="M24 20.993V24H0v-2.996A14.977 14.977 0 0112.004 15c4.904 0 9.26 2.354 11.996 5.993zM16.002 8.999a4 4 0 11-8 0 4 4 0 018 0z" />
    </svg>
  </span>
  <div class="space-y-1">
    <div class="text-sm font-medium text-gray-900 dark:text-white">Anonymous</div>
    <a href="#" class="text-xs text-gray-500 dark:text-gray-200 font-mono">
      {shortPeerID(peerID)}
    </a>
  </div>
</div>