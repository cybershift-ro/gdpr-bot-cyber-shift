[![Banner](https://readme-typing-svg.demolab.com?font=Source+Code+Pro&size=24&duration=4000&pause=100&color=213E68&center=true&multiline=true&width=500&height=100&lines=GDPR+Bot;Cyber+Shift+%26+Privacy+Magistrus)](https://git.io/typing-svg)

# Intro
  
GDPR Bot is a small web scraper that keeps you updated on the latest sanctions announced by [The National Authority for the Supervision of the Processing of Personal Data](https://www.dataprotection.ro/?page=allnews).  
  
The scapper crawls [the news page](https://www.dataprotection.ro/?page=allnews) once 30 minutes and builds a complete feed available in 3 formats: [RSS][2], [ATOM][3] and [JSON][4].
  

It is written in Go programming language and compiles down to a single binary that can be executed on any OS.

# Motivation
  
We want to show a wider audience that there are many vulnerable companies in Romania due to an incorrect GDPR implementation.

The website [ANSPDCP Dataprotection Romania][1] offers, unfortunately, a short RSS feed, consisting only of the last 5 announcements. These announcements are not necessarily related to the sanctions imposed on companies that do not comply with the GDPR policy.

GDPR Bot offers a complete list of all sanctions announced by [ANSPDCP Dataprotection Romania][1].

[Cyber Shift](https://cybershift.dev) together with the [Privacy Magistrus](https://gdprmag.com/) team is developing this program transparently.

# Demo
  
  A live version is available at [gdpr.cybershift.dev](https://gdpr.cybershift.dev/rss).

# Features
  - Fast & flexible crawler with the help of [Colly](https://github.com/gocolly/colly)
  - Web server built-in
  - Sanctions feed available in 3 formats: [RSS][2], [ATOM][3] and [JSON][4].
  - Structured information about [sanctions found][5]
  - Integrated backend with front-end - [Pocket Base](https://pocketbase.io/)
  - Compatible with any RSS reader
  - Easy to self-host

# Upcoming features
  - Automatically extract company name and sanction total value from the news article
  - Config file

[1]: <www.dataprotection.ro/> "ANSPDCP Dataprotection Romania"
[2]: <https://gdpr.cybershift.dev/rss>  "RSS Feed of GDPR Bot"
[3]: <https://gdpr.cybershift.dev/atom> "ATOM Feed of GDPR Bot"
[4]: <https://gdpr.cybershift.dev/json> "JSON Feed of GDPR Bot"
[5]: <https://gdpr.cybershift.dev/sanctions> "JSON Data of Fines found by GDPR Bot"
