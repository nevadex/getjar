# getjar
### a command line utility to download or compile minecraft server jars

---

Why did I make this? MCJarFinder was feeling limited and it was getting annoying having to figure out how to wget my way into downloading all the jars and which websites have adwalls so this is now a thing

---

## currently supported:

- vanilla (version_manifest)
- spigot (spigotmc buildtools)
- craftbukkit (spigotmc buildtools)
- paper (downloads api)
- folia (downloads api)

### in development:

- mohist (downloads api)
- fabric (downloads api)
- catserver (website hotlinks cuz only 3 versions lol)
- purpur (downloads api)
- also going to make this work as a library available as "github.com/nevadex/getjar/lib"

award-winner for the shiddiest download platform:  
forge - adfocus links, hate it and plus the maven is locked

> jar options are based on what is available on serverjars.com but getjar is designed to operate without a middleman website or api, unlike serverjars.com

---

## limitations:

- any type of jar without an indexable source (a list of versions or a list of the latest versions in all channels) will probably not make it into this unless someone makes a pull request (soup ftw)
- no version lists (inconsistency between minecraft version and project version, some don't have lists, etc)