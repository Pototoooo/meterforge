# Documentation requirements

The `doc-decorator` linter rule (warning) requires a `/** comment */` or `@doc` on:

- All named models, enums, and unions.
- All model properties, except `_` and `contentType`.

Operations must have both `@operationId` (kebab-case) and `@summary`.

Use `#suppress "@meterforge/api-spec-aip/doc-decorator" "<reason>"` only for shared base model spreads where documenting each inherited property individually would be redundant.
