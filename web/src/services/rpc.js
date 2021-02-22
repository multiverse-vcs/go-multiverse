// sequence is the unique request id
let sequence = 0

// body returns the rpc request body string
function body(method, params) {
  return JSON.stringify({
    id: sequence++,
    method: method,
    params: params,
  })
}

// request invokes a remote method
function request(method, params) {
  return fetch('http://localhost:9001/', {
    cache: 'no-cache',
    method: 'POST',
    body: body(method, params),
    headers: {
      'Content-Type': 'application/json'
    },
  }).then(res => res.json())
}

export default {
  Author: {
    Search: (peerID) => request('Author.Search', [{peerID}]),
    Self: () => request('Author.Self', []),
  },
  Repo: {
    Fetch: (remote) => request('Repo.Fetch', [{remote}]),
    List: () => request('Repo.List', []),
  }
}