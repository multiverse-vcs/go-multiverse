const RPC_URL = "http://localhost:9001/_jsonRPC_";
const HTTP_URL = "http://localhost:2020";

let sequence = 0;

function rpc(method, params) {
  return fetch(RPC_URL, {
    cache: "no-cache",
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id: sequence++,
      method: method,
      params: params,
    }),
  });
}

export default {
  Self: async function () {
    const res = await rpc("Author.Self", []);
    return await res.json();
  },

  SearchAuthor: async function (peerID) {
    const res = await rpc("Author.Search", [{ peerID }]);
    return await res.json();
  },

  ListRepos: async function () {
    const res = await rpc("Repo.List", []);
    return await res.json();
  },

  FetchRepo: async function (remote) {
    const res = await fetch(`${HTTP_URL}/${remote}`);
    return await res.json();
  },

  FetchFile: async function (remote, branch, path) {
    const res = await fetch(`${HTTP_URL}/${remote}/${branch}/${path}?highlight=monokai`);
    const type = res.headers.get("Content-Type");

    if (type === "application/json") {
      return await res.json();
    }

    return await res.text();
  },
};
