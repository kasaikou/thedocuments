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
  color       = "B6D078"
  description = "Create application, service or tools"
}

resource "github_issue_label" "creation" {
  for_each   = local.github_repositories
  repository = each.value

  name        = "creation"
  color       = "D18C07"
  description = "Output something without articles"
}

resource "github_issue_label" "other" {
  for_each   = local.github_repositories
  repository = each.value

  name        = "other"
  color       = "BFDADC"
  description = "Other issues (not labeled)"
}

