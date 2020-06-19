This is the source for programs available on [chaos.projectdiscovery.io](http://chaos.projectdiscovery.io/). Please send pull-request of public bug bounty programs that you want to include in our public list with recon data. 

We are currently accepting in JSON format, an example is below:

```json
{
   "name":"HackerOne",
   "url":"https://hackerone.com/security",
   "bounty":true,
   "domains":[
      "hackerone.com",
      "hackerone.net",
      "hacker101.com",
      "hackerone-ext-content.com"
   ]
}
```

**Notes:** 
1. Use JSON validators (e.g. https://jsonlint.com) to validate the modfied `bugbounty-list.json` file when sending pull requests.
2. Scope file only accept root domain as input the `domains` field, so please do not add subdomains in `domain` field. 

Thanks again for your contribution and keeping the community vibrant. :heart:

-------

If you want to remove any program from the list, please contact us at contact@projectdiscovery.io.
