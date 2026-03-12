# Features

**NOTE:** This is an experimental library. Stability is not guaranteed.

Features is a Go library for building flexible, schema-validated data models. It pairs immutable persistent maps with JSON Schema validation, so you get the freedom of dynamic data structures without giving up the safety of a well-defined schema.

## Core Concepts

### Features

A feature is the central abstraction: a key-value map with a schema attached. You can set and retrieve typed values — strings, numbers, integers — and validate the entire structure against its schema at any point. Because features are backed by persistent data structures, every modification produces a new version without copying the entire map.

### Schemas and Migrations

A schema is not written as a single definition. Instead, it is the sum of its migrations — an ordered sequence of operations that each add or remove fields. When you need to validate data, the library reduces all migrations into a single [JSON Schema](https://json-schema.org/) and checks your feature against it.

This design means that evolving your data model is a first-class concern, not an afterthought. Adding a required field with a default value, removing an obsolete one, or restructuring your schema over time is expressed as a series of small, composable steps.

#### Operations

Migrations are composed of operations:

- **AddField** — introduces a new field with a name, type, and optional default value.
- **RemoveField** — drops a field from the schema.

### Forms

Forms bridge the gap between raw user input and validated features. A form wraps a feature and can populate it directly from URL query parameters or POST data, coercing string values into the correct types based on the schema. After populating the form, a single `Validate` call checks the entire feature against its schema and returns any errors.

### Keys

Features supports flexible key generation for storage. You can compose keys from literal strings, feature property values, or combinations of both. This makes it straightforward to build partition keys, sort keys, or any other indexing scheme your storage layer requires.

### Reconciliation

An early-stage three-way merge system for resolving concurrent edits to the same feature. It compares an incoming change against a shared base and the current head, using structural diffs to detect and surface conflicts.

## Installation

```
go get github.com/mamaar/features
```
