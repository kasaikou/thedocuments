locals {
  github_repositories = toset([
    "thedocuments",
    "awesomedocuments"
  ])
}

resource "github_issue_label" "documentation" {
  for_each   = local.github_repositories
  repository = each.value

  name        = "documentation"
  color       = "0075CA"
  description = "Improvements or additions to documentation"
}

resource "github_issue_label" "software" {
  for_each   = local.github_repositories
  repository = each.value

  name        = "software"
  color       = "D18C07"
  description = "Create application, service or tools"
}

resource "github_issue_label" "other" {
  for_each   = local.github_repositories
  repository = each.value

  name        = "other"
  color       = "BFDADC"
  description = "Other issues (not labeled)"
}

