# How to Contribute

We'd love to accept your patches and contributions to this project. There are
just a few small guidelines you need to follow.

## Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License
Agreement. You (or your employer) retain the copyright to your contribution,
this simply gives us permission to use and redistribute your contributions as
part of the project. Head over to <https://cla.developers.google.com/> to see
your current agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.

## Code of Conduct

Participation in this project comes under the [Contributor Covenant Code of Conduct](code-of-conduct.md)

## Submitting code via Pull Requests

*Thank you* for considering submitting code to Agones!

- We follow the [GitHub Pull Request Model](https://help.github.com/articles/about-pull-requests/) for
  all contributions.
- For large bodies of work, we recommend creating an issue and labelling it
  "[kind/design](https://github.com/googleprivate/agones/issues?q=is%3Aissue+is%3Aopen+label%3Akind%2Fdesign)"
  outlining the feature that you wish to build, and describing how it will be implemented. This gives a chance
  for review to happen early, and ensures no wasted effort occurs.
- For new features, documentation *must* be included. Review the [Documentation Editing and Contribution](https://agones.dev/site/docs/contribute/)
  guide for details.
- It is strongly recommended that new API design follows the [Google AIPs](https://google.aip.dev/) design guidelines.  
- All submissions, including submissions by project members, will require review before being merged.
- Once review has occurred, please rebase your PR down to a single commit. This will ensure a nice clean Git history.
- If you are unable to access build errors from your PR, make sure that you have joined the [agones-discuss mailing list](https://groups.google.com/forum/#!forum/agones-discuss).
- Please follow the code formatting instructions below.

## Formatting

When submitting pull requests, make sure to do the following:

- Format all Terraform code with `terraform fmt`. 
- Format all Go code with [gofmt](https://golang.org/cmd/gofmt/). Many people
  use [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) which
  fixes import statements and formats code in the same style of `gofmt`.
- C++ code should follow the [Google C++ Style
  Guide](https://google.github.io/styleguide/cppguide.html), which can be
  applied automatically using the
- Remove trailing whitespace. Many editors will do this automatically.
- Ensure any new files have [a trailing newline](https://stackoverflow.com/questions/5813311/no-newline-at-end-of-file)

## Feature Stages

Often, new features will need to go through experimental stages so that we can gather feedback and adjust as necessary.

You can see this project's [feature stage documentation](https://agones.dev/site/docs/guides/feature-stages/) on the Agones
website.

If you are working on a new feature, you may need to take feature stages into account. This should be discussed on a
 design ticket prior to commencement of work. 

