# Database Migrations in MeterForge

MeterForge uses [Ent](https://entgo.io) to define its database schema and Atlas to generate versioned migrations. Versioned migrations are the only supported way to initialize and upgrade the database.

## AutoMigrate

MeterForge can apply or wait for versioned migrations during startup. This behavior is configured through `postgres.autoMigrate`.

- `migration` applies the embedded versioned migration history.
- `migration-job` waits until a separate migration job has brought the database to the version required by the running binary.
- `false` disables startup migration handling for installations that manage migrations externally.

Ent automigration is deprecated and no longer supported. The former `postgres.autoMigrate: ent` value is rejected with instructions for adopting the database. The `postgres.autoMigrate` setting itself remains supported for versioned migrations through `migration`, `migration-job`, and `false`.

## Upgrading from Ent-managed schema migration

Older MeterForge installations could use `postgres.autoMigrate: ent`. In that mode, MeterForge used Ent's additive schema migration to synchronize database tables at startup instead of applying the versioned Atlas migration history.

The Atlas migrations also contain data backfills and database objects that cannot be expressed by Ent's table schema, including functions, triggers, views, and required singleton rows. Continuing to support both paths would allow Ent-managed databases to diverge from databases upgraded through Atlas. Going forward, MeterForge supports only the Atlas-based versioned migrations.

The `adopt-ent` command is the one-time bridge for installations upgrading from `postgres.autoMigrate: ent`. It only accepts a non-empty MeterForge database without `schema_om`, which identifies a legacy Ent-managed database. Empty and already-versioned databases must use the normal migration command.

Stop all MeterForge servers and workers before starting the upgrade. Keeping MeterForge stopped until both commands complete prevents application processes from reading or modifying a partially adopted schema.

First, adopt the Ent-managed database into versioned migration ownership:

```bash
meterforge-jobs migrate adopt-ent
```

This command:

1. Applies the additive Ent schema behavior frozen at MeterForge commit `12ab7b082035f2f93972c7f98973c5502107c157`.
2. Applies the frozen reconciliation scripts for data backfills and database objects not represented by Ent.
3. Verifies the reconciled state.
4. Records migration baseline `20260709134422` in `schema_om`.

It deliberately stops at that baseline. Next, run the unchanged normal migration command:

```bash
meterforge-jobs migrate
```

The normal command applies every migration after `20260709134422` up to the version embedded in the target MeterForge release. This procedure therefore works regardless of the target version: `adopt-ent` always establishes the same baseline, and the normal migration history performs the remaining upgrade.

The adoption command refuses empty databases, already-versioned databases, and unversioned databases that are not recognizable as MeterForge. Take a database backup before adoption. Start MeterForge again only after the normal migration command succeeds.

## Generating Migrations

MeterForge uses [Atlas](https://atlasgo.io/) to generate versioned migrations from changes in the Ent schema.

After changing the schema and running `go generate` you can create a new migration diff via `atlas migrate --env local diff <migration-name>`, the generated migration files will be placed in the `migrations` directory.

## Data Migrations

Data migrations can be written alongside the schema migrations code.

### Testing Data Migrations

To test your data migrations in hypothetical scenarios, you can use the `tools/migrate` package to write assertions on how data changes after applying the migration. Examples can be found in the package.

## Running Migrations

The recommended way to run normal versioned migrations is the bundled migration job:

```bash
meterforge-jobs migrate
```

Legacy Ent-managed databases must first run `meterforge-jobs migrate adopt-ent` as described above.

For already-versioned databases, the `golang-migrate` CLI can also apply the migration directory directly:

```bash
migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&x-migrations-table=schema_om" -path ./tools/migrate/migrations up
```

MeterForge records migration state in `schema_om`.
