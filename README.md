# Leagueify

Leagueify is an open source league management software that allows you to develop competitive and fun leagues for teams of all ages.

This server is written in [Go][go-website] using the [Echo][echo-website] framework.

## Getting Started

Using the provided development Docker compose example, Leagueify is quick and easy to get running with intuitive make commands.

```bash
# Start Leagueify in development mode
make start-dev

# Start Leagueify in detached development mode
make start-dev-detached

# Stop Leagueify and remove docker images
make clean-dev
```

After starting Leagueify with `make start-dev` you should see the Leagueify banner in the terminal, informing that the system is ready to be used.
Upon saving changes, [air][air-github] will automatically reload the system, Leagueify will be ready once the banner is again shown.

[air-github]: https://github.com/air-verse/air
[echo-website]: https://echo.labstack.com
[go-website]: https://go.dev
