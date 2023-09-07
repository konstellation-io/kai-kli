## [v0.2.0](https://github.com/konstellation-io/kai-kli/tag/v0.2.0) (2023-09-07)

### Fixed

* [[bc173](https://github.com/konstellation-io/kai-kli/commit/bc173a6aa3f82f222dff5e8aafa31c9146ff368d)] fix: validate auth url protocol schema

 (Juan Jesús Padrón Hernández, 2023-09-05T08:14:22)

* [[8afaa](https://github.com/konstellation-io/kai-kli/commit/8afaa2eb90d740ddf52df97c98774381b6f32ed4)] fix: fix url validation on login command

 (David Fuentes, 2023-08-31T07:39:31)

* [[1ed35](https://github.com/konstellation-io/kai-kli/commit/1ed35566f5b3bd9eaeed5ee8022bd47d23015dd3)] fix: rename kai related concepts

 (Juan Jesús Padrón Hernández, 2023-07-05T18:50:38)

* [[90224](https://github.com/konstellation-io/kai-kli/commit/902242de9d9e64231ca05bf882a47be70a8da9b3)] fix: rename runtime to product

fix: rename runtime to product (Juan Jesús Padrón Hernández, 2023-05-29T13:45:35)

* [[72c4c](https://github.com/konstellation-io/kai-kli/commit/72c4c5f2ae255542b44f2b5641a86cfb717f38c4)] fix: rename update version configuration to correct grahpql mutation ([#3](https://github.com/konstellation-io/kai-kli/issues/3))

&lt;!--- Provide a general summary of your changes in the Title above --&gt;

## Description
&lt;!--- Describe your changes in detail --&gt;
updateVersionConfiguration renaming to correct graphql mutation
updated help doc string for update version config cmd

## Motivation and Context
&lt;!--- Why is this change required? What problem does it solve? --&gt;
&lt;!--- If it fixes an open issue, please link to the issue here. --&gt;
updateVersionConfiguration graphql is not the correct mutation name

## Types of changes
&lt;!--- What types of changes does your code introduce? Put an `x` in all the boxes that apply: --&gt;
- [X] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Other changes (ci configuration, documentation or any other kind of changes)

## Checklist:
&lt;!--- Go over all the following points, and put an `x` in all the boxes that apply. --&gt;
&lt;!--- If you&apos;re unsure about any of these, don&apos;t hesitate to ask. We&apos;re here to help! --&gt;
- [ ] I have created tests for my code changes, and the tests are passing.
- [ ] I have executed the pre-commit hooks locally.
- [ ] My change requires a change to the documentation (create a new issue if the documentation has not been updated).
- [ ] I have updated the documentation accordingly.
 (Alex Rodriguez Fernandez, 2023-05-29T10:45:20)

* [[249ab](https://github.com/konstellation-io/kai-kli/commit/249abb0282e3fb79ea42e8e4c9535f75faf09298)] fix(graphql): updte the graphql library to the latest version

 (David Fuentes, 2023-05-29T06:46:57)

* [[98879](https://github.com/konstellation-io/kai-kli/commit/98879294db5b72cd19e4a74b4f87e27a2e115829)] fix: typo in test variable ([#59](https://github.com/konstellation-io/kai-kli/issues/59))

Co-authored-by: javier.aguilera &lt;javier.aguilera@intelygenz.com&gt; (kafkaphoenix, 2023-03-08T08:30:48)

* [[d5291](https://github.com/konstellation-io/kai-kli/commit/d52913f07bec14f83d65af18432fd292d6e6bb7a)] fix: Sonarcloud integration

 (David Fuentes, 2023-02-15T09:20:51)

### Added

* [[4214f](https://github.com/konstellation-io/kai-kli/commit/4214f16ee722a32d0a0dea29cf8bcff82af13e84)] Merge pull request [#3](https://github.com/konstellation-io/kai-kli/issues/3)2 from konstellation-io/feature/add-scoop-to-goreleaser

feat: add scoop to goreleaser (Ángel Luis Piquero Coloma, 2023-08-29T12:35:06)

* [[88b9c](https://github.com/konstellation-io/kai-kli/commit/88b9c24d058e5c1b62eaed7c7774ada378ed740b)] feat: Remove credentials from authentication

 (David Fuentes, 2023-08-28T09:34:18)

* [[ee4c8](https://github.com/konstellation-io/kai-kli/commit/ee4c82bc224bb07a71011b90639f2db1bcd25039)] feat: process-registry register command

 (David Fuentes, 2023-08-24T13:11:01)

* [[a85d2](https://github.com/konstellation-io/kai-kli/commit/a85d241918cc51a42372b9cdd44e53c6f44d1b3e)] feat: add workflow and process management ccommands, add configuration management command, fix login for IdPs

 (David Fuentes, 2023-08-11T07:34:49)

* [[1baf4](https://github.com/konstellation-io/kai-kli/commit/1baf4ef932406d61a0d27c919bac9a818e172cef)] feat: logout command

 (David Fuentes, 2023-07-21T10:03:46)

* [[98d35](https://github.com/konstellation-io/kai-kli/commit/98d35c16fcdb7c8d5337e0e6bf5e7081b5ce05c5)] feat: add login command

 (David Fuentes, 2023-07-18T13:41:56)

* [[bb6a2](https://github.com/konstellation-io/kai-kli/commit/bb6a2736bb39b5dfdb3b3316f90ef25dbca000c5)] feat: Add remove server command 

 (David Fuentes, 2023-07-13T14:06:45)

* [[00b62](https://github.com/konstellation-io/kai-kli/commit/00b62f8fefcc4bc3e3df55f6b2262e69cddaabbf)] feat: add setup step as a PersistentPreRun

 (David Fuentes, 2023-07-11T16:10:28)

* [[4b616](https://github.com/konstellation-io/kai-kli/commit/4b6162a46d4be17a1b1c89ac9b210aecd9e95410)] feat: add multi runtime support ([#58](https://github.com/konstellation-io/kai-kli/issues/58))

Co-authored-by: David Fuentes &lt;david.fuentes@intelygenz.com&gt; (Juan Jesús Padrón Hernández, 2023-02-27T17:01:17)

