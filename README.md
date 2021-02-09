# Public BugBounty Programs

[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/projectdiscovery/public-bugbounty-programs/issues)
[![Follow on Twitter](https://img.shields.io/twitter/follow/pdchaos.svg?logo=twitter)](https://twitter.com/pdchaos)
[![Chat on Discord](https://img.shields.io/discord/695645237418131507.svg?logo=discord)](https://discord.gg/KECAGdH)

This is a source of public programs listed on [chaos.projectdiscovery.io](https://chaos.projectdiscovery.io/). Please send pull-request of public bug bounty programs that you want to include in our public list with recon data. 

We are currently accepting in JSON format, an example is below:

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


💬 Discussion
-----

Have questions / doubts / ideas to discuss? feel free to open a discussion using [Github discussions](https://github.com/projectdiscovery/public-bugbounty-programs/discussions) board.

👨‍💻 Community
-----

You are welcomed to join our [Discord Community](https://discord.gg/KECAGdH). You can also follow us on [Twitter](https://twitter.com/pdchaos) to keep up with everything related to projectdiscovery, got question? please reach out to us at chaos@projectdiscovery.io

📋 Notes
-----
- Only domain name values are accepted in the `domains` field.
- We do not support wildcard input like `*.tld` or `*.tld.*`.
- **domains** field includes TLD names associated with the target program, not based on scope of the program.
- Subdomains are populated using **Passive API** (chaos dataset).


📌 Reference
-----

- https://github.com/arkadiyt/bounty-targets-data
- https://github.com/disclose/diodb/blob/master/program-list.json
- https://firebounty.com

Thanks again for your contribution and keeping the community vibrant. :heart: