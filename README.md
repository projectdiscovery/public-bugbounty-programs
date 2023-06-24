<h1 align="center"> Public BugBounty Programs </h1>


<h4 align="center"> Community curated list of public bug bounty and responsible disclosure programs. </h4>

<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://github.com/projectdiscovery/public-bugbounty-programs/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://twitter.com/pdiscoveryio"><img src="https://img.shields.io/twitter/follow/pdiscoveryio.svg?logo=twitter"></a>
<a href="https://discord.gg/projectdiscovery"><img src="https://img.shields.io/discord/695645237418131507.svg?logo=discord"></a>
</p>

The [chaos-bugbounty-list.json](chaos-bugbounty-list.json) file serves as the central management system for the public bug bounty programs displayed on [chaos.projectdiscovery.io](https://chaos.projectdiscovery.io/). We welcome your contributions to this list. If there are specific programs for which you'd like to see reconnaissance data, please submit a pull request.

We are currently accepting submissions in JSON format. Here's an example of the structure we require:

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

Your contributions will help us to continually improve and expand the range of public bug bounty programs we feature.


💬 Discussions
-----

For any inquiries, suggestions, or topics you'd like to discuss, we encourage you to initiate a "Discussion" using our [GitHub Discussions](https://github.com/projectdiscovery/public-bugbounty-programs/discussions) platform.

👨‍💻 Community
-----

We invite you to join our [Discord Community](https://discord.gg/projectdiscovery) for more interactive discussions.  
Stay updated with our latest news and activities by following [ProjectDiscovery](https://twitter.com/pdiscoveryio) on Twitter.  
For direct communication, feel free to reach us at [contact@projectdiscovery.io](mailto:contact@projectdiscovery.io).

📋 Guidelines
-----
- Please note that only domain name values are accepted in the `domains` field.
- We do not support wildcard inputs such as `*.tld` or `*.tld.*`.
- The **domains** field should include TLD names associated with the target program, not necessarily based on the scope of the program.
- Subdomains are populated using our dataset **Passive API**

📌 References
-----

- https://github.com/arkadiyt/bounty-targets-data
- https://github.com/disclose/diodb/blob/master/program-list.json
- https://firebounty.com

We greatly appreciate your contributions and your efforts in keeping our community dynamic and engaging. :heart: