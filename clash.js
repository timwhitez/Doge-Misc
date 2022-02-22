module.exports.parse = ({ content, name, url }, { yaml, axios, notify }) => {
  proxies = [];
  for (let proxy of content.proxies) {
    if (proxy.server === undefined) continue;

      proxies.push(proxy.name);

  }
  if (proxies.length > 0) {
    content['proxy-groups'].push({
      'name': '秒切',
      'type': 'load-balance',
      'proxies': proxies,
      'url': 'http://www.gstatic.com/generate_204',
      'strategy': 'round-robin',
      'interval': "1"
    });
    content['proxy-groups'][0].proxies.unshift("秒切");
  }
  return content;
}
