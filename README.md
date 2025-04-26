# getjar
### a command line utility to download or compile minecraft server jars

---

why did I make this? mcjarfinder was feeling limited and it was getting annoying having to figure out how to wget my way into downloading all the jars and which websites have adwalls so this is now a thing

---

## install:

if go is installed: `go install github.com/nevadex/getjar@latest`

if not: download binary from the "Releases" page and add to path

## usage:

`getjar [type] [flags]`

eg. `getjar vanilla -v 1.13.2 -f aquarium.jar`

use `getjar -h` or `getjar [type] -h` for more options

---

## features:

- download server jars directly
- download supported version lists directly
- select specific builds/internal versions
- access a wide range of popular server options
- available as a library for other applications (github.com/nevadex/getjar/ops)

---

## currently supported:

- **vanilla** (version_manifest)
- **spigot** (spigotmc buildtools)
- **bukkit** (spigotmc buildtools)
- **paper** (papermc api)
- **folia** (papermc api)
- **mohist** (mohistmc api)
- **banner** (mohistmc api)
- **fabric** (fabric meta api)
- **catserver** (jenkins api)
- **purpur** (purpurmc api)
- **forge** (html parsing, later replace with maven)
- **neoforge** (maven)


### in development:

- magma (magmafoundation downloads api)
- sponge (spongepowered downloads api)
- limbo (jenkins api)
- pufferfish (jenkins api)
- send any other wanted server types as an issue on this repo (or write it yourself and make a pull request)

> jar options are based on what ~~is~~ was available on serverjars.com but getjar is designed to operate without a middleman website or api, unlike serverjars.com

---

## limitations:

- any type of jar without an indexable source (a list of versions or a list of the latest versions in all channels) will probably not make it into this