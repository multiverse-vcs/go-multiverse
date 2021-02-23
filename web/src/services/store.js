import { readable, writable, derived } from "svelte/store";
import api from "./api";

export const self = readable(null, (set) => {
	api.Self().then((res) => set(res.result));
	return () => {};
});

export const remote = writable(null);
export const branch = writable(null);
export const path = writable("");

export const repo = derived(remote, ($remote, set) => {
	if (!$remote) return;
	api.FetchRepo($remote).then((res) => set(res));
});

export const file = derived([remote, branch, path], ([$remote, $branch, $path], set) => {
		if (!$remote || !$branch) return;
		api.FetchFile($remote, $branch, $path).then((res) => set(res));
	}
);
