import { readable, writable, derived } from 'svelte/store'
import rpc from './rpc'

// peerID is the current user peer ID
export const peerID = readable(null, set => {
  rpc.Author.Self().then(res => set(res.result.peerID))
  return () => {}
})

// authorID is the current author peer ID
export const authorID = writable(null, set => {
  rpc.Author.Self().then(res => set(res.result.peerID))
  return () => {}
})

// author is the current author object
export const author = derived(authorID, ($authorID, set) => {
  if (!$authorID) return () => {}
  rpc.Author.Search($authorID).then(res => set(res.result.author))
  return () => {}
})

// repoName is current repository name
export const repoName = writable(null)

// repo is the current repository object
export const repo = derived([authorID, repoName], ([$authorID, $repoName], set) => {
  if (!$authorID || !$repoName) return () => {}
  rpc.Repo.Fetch(`${$authorID}/${$repoName}`).then(res => set(res.result.repository))
  return () => {}
})