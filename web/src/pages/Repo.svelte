<script>
  import { onDestroy } from 'svelte'
  import { repoName, repo } from '../services/store'

  let branches = []

  const unsubscribe = repo.subscribe(value => {
    if (!value) return

    branches = Object.keys(value.branches).sort()
  })

  onDestroy(unsubscribe)
</script>

<h1>{$repoName}</h1>
{#each branches as branch}
  <h2>{branch}</h2>
{/each}