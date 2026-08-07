# Upstream and project scope

## Source project

MeterForge is derived from the open-source [OpenMeter](https://github.com/openmeterio/openmeter) codebase.

The imported history is retained in this repository. The local commit immediately before the MeterForge namespace migration is:

```text
0ee5b9817 feat(charges): support custom currency credits-only charges (#4769)
```

The repository-level MeterForge migration is recorded as:

```text
7d26359d3 Rebrand project as MeterForge
```

These identifiers describe this repository's rewritten history and are not asserted to be identical to commit IDs in the upstream repository.

## MeterForge-specific scope

The current repository focuses on:

- migrating the application and code namespace to MeterForge;
- providing a self-hosted Console oriented toward local product exploration;
- Chinese UI and documentation for the primary metering and billing workflow;
- AI and Agent usage scenarios, especially Token-based metering;
- local Docker Compose deployment and reproducible study workflows.

Generated SDKs and API artifacts are produced from the repository's TypeSpec and OpenAPI sources. Package names visible in source do not imply that a package has been published to a public registry.

## Attribution and license

The upstream project is licensed under the Apache License 2.0. MeterForge retains that license in [LICENSE](LICENSE). Copyright and attribution notices in source files and bundled third-party components remain subject to their respective terms.
