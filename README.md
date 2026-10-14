# What's that?

This is a study project I assembled while following the book "Let's Go Further" by
Alex Edwards, with some minor tweaks for my environment. It is a JSON API for
retrieving and managing movie records.

The project is written test-first: every feature starts as a failing test.

## How to start

I don't expect anyone but myself to run this app, but if you want to, here's how:

1. Install Go, Docker and the `migrate` CLI:

   ```sh
   brew install golang-migrate
   ```

2. Prepare the environment and start PostgreSQL: `make prepare-env db-up`.
3. Apply the migrations: `make migrate-up`.
4. Run the app: `make run`.

That's it!

`make prepare-env` copies `.env.example` to `.env`; edit it if you want a
different database user, password or name. The same values are picked up by
`docker-compose.yml`, so the container and the app always agree.
