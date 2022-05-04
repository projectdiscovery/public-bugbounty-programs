# Public BugBounty Programs

[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/projectdiscovery/public-bugbounty-programs/issues)
[![Follow on Twitter](https://img.shields.io/twitter/follow/pdchaos.svg?logo=twitter)](https://twitter.com/pdchaos)
[![Chat on Discord](https://img.shields.io/discord/695645237418131507.svg?logo=discord)](https://discord.gg/projectdiscovery)

[This](chaos-bugbounty-list.json) JSON file controls the public bug bounty programs listed on [chaos.projectdiscovery.io](https://chaos.projectdiscovery.io/). Please create a pull-request with the programs for which you'd like to see recon data. 

We are currently accepting entries in JSON format. See an example below:

```json
{
   "name":"HackerOne",
   "url":"https://hackerone.com/security",
   "bounty": true,
   "swag": true,
   "domains":[
      "hackerone.com",
      "hackerone.net",
      "hacker101.com",
      "hackerone-ext-content.com"
   ]
}
```


💬 Discussions
-----

If you have any questions/doubts/ideas to discuss, please create a "Discussion" using the [GitHub Discussions](https://github.com/projectdiscovery/public-bugbounty-programs/discussions) board.

👨‍💻 Community
-----

Join our [Discord Community](https://discord.gg/projectdiscovery).  
Follow [@PDChaos](https://twitter.com/pdchaos) and [PDiscoveryIO](https://twitter.com/pdiscoveryio) on Twitter.  
You can also contact us at [chaos@projectdiscovery.io](mailto:chaos@projectdiscovery.io).

📋 Notes
-----
- Only domain name values are accepted in the `domains` field.
- We do not support wildcard input like `*.tld` or `*.tld.*`.
- **domains** field includes TLD names associated with the target program, not based on scope of the program.
- Subdomains are populated using **Passive API** (chaos dataset).


📌 References
-----

- https://github.com/arkadiyt/bounty-targets-data
- https://github.com/disclose/diodb/blob/master/program-list.json
- https://firebounty.com

Thank you for your contribution and for keeping the community vibrant. :heart: