TLDR: Versioning the documentation adds more work than its worth. Maybe reconsider when releasing a lts release.


> You release v1, and start immediately working on v2 (including its
> docs). In this case, the current version is v2, which is in the ./docs
> folder, while the latest version is v1, which is the version hosted at
> example.com/docs and is browsed by most of your users.

> You release v1,
> and will maintain it for some time before thinking about v2. In this
> case, the current version and latest version will both be point to v1,
> since the v2 docs doesn't even exist yet!

From the [Docusaurus docs][docu-ver]

#### [Method 1 Prototype branch][gh-method-1] (default docs are in `/website/versioned_docs/version-***/`)
#### [Method 2 Prototype branch][gh-method-2](default docs are in `./docs/`)


# Potential Universe Documentation Issues

The packages downloaded by `dagger project update` might get out-of-sync 
with the documentation. I.e universe package `foo 1.2.3` is downloaded with `dagger 4.5.6`
but `foo 0.1.0` documentation is shown under the `dagger 4.5.6` version.


# Other Notes 
  - The dropdown menu has the same font & background color on the default theme. 
  - Most `make` targets acting on the docs will break. (`make web` still works!).
  - **lots** of duplicate pages. A typo would have to be changed in every version.
  - Unable to find a way to change the path of `versioned_docs` and `versioned_sidebars`
  - Warning Banner when viewing old/unreleased docs.

<!-- links -->
[docu-ver]: https://docusaurus.io/docs/versioning

[gh-method-1]: https://github.com/KGB33/dagger/tree/futureDocs-proto
[gh-method-2]: https://github.com/KGB33/dagger/tree/doc-versioning-proto
