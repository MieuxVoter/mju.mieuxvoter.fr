# Majority Judgment Urn Companion

[![Discord Chat https://discord.gg/k9YRuZPSZs](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/k9YRuZPSZs)


Try it out at https://mju.mieuxvoter.fr


## Purpose

Provide a Digital Asset Links file so that Android can verify our app.

https://github.com/MieuxVoter/majority-judgment-offline-urn-android

We might also redirect users to places where they can download the app.

> That's pretty much it.  This website could just be static files.


## Developers — Run locally

```shell
make build && build/mju
```

or, using Docker:

```shell
make start
```

Then, visit http://localhost:8034


### Watch for changes and reload

During development, for fast iterations, it's handy to reload the webserver when a file changes.
I'm using `entr` for this, it works like a charm:

```bash
find src | entr -r go run src/main.go
```
