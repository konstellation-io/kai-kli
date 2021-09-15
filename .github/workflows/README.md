# Workflows

Workflow uses actions defined in [this repository](https://github.com/intelygenz/action-product-version-tags) and extra
info can be found in the repository README.md

These are the workflows defined in the repository, in logical order of execution:

1. Quality
2. Pre Release
3. Release
4. Go Releaser

## Quality

Run quality actions (list, tests, coverage)

Conditions:

- On push for branches not starting with `v` or tags not starting with `v`

## Pre Release

Generates a new alpha release tag.

Conditions:

- Only in `main` branch and only when quality workflow was triggered

## Release

Generates a new release tag and branch. The tag is calculated taking in account the last pre-release tag.

Conditions:

- Manual run

## Go Releaser

Generates a new release of the code

Conditions:

- On tags `v*`
