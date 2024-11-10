# Features

## Flexible data models

Flexible data models allow you to define and manipulate data structures dynamically. 
This module provides tools to create, modify, and query these models efficiently. 
The underlying CHAMP structure ensures that operations on these models are both fast and memory-efficient.

## Schema validation

Flexible data models are great, but we also need to ensure that the data we store is valid.
On top of the flexible maps, we provide what we call `features` which are
a flexible map with a schema attached to it. This schema is used to validate the data.

## Schema

A schema is the sum of its migrations. Under the hood, a schema is defined as a [JSON schema](https://json-schema.org/).

### Migrations

A migration is a set of operations that can be applied to a schema to transform it into another schema.
For validation the migrations are applied in order to build the final schema.

### Operations

There are various types of operations:
- `AddField`
- `RemoveField`

## Forms

A form is a visual representation of a feature. It can contain any subset of fields from the feature.
The form has a validate method that will validate the entire feature against its schema.