resource "google_service_account" "cloudbuild-sa" {
  project      = var.project
  account_id   = "cloudbuild-cicd"
  display_name = "Cloud Build - CI/CD service account"
}

resource "google_project_iam_member" "cloudbuild-sa-cloudbuild" {
  project = var.project
  role    = "roles/cloudbuild.builds.builder"
  member  = "serviceAccount:${google_service_account.cloudbuild-sa.email}"
}

resource "google_project_iam_member" "cloudbuild-sa-cloudbuild-roles" {
  project = var.project
  for_each = toset([
    "roles/serviceusage.serviceUsageAdmin",
    "roles/clouddeploy.operator",
    "roles/cloudbuild.builds.builder",
    "roles/container.admin",
    "roles/storage.admin",
    "roles/iam.serviceAccountUser"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.cloudbuild-sa.email}"
}

resource "google_project_iam_member" "clouddeploy-admin" {
  project = var.project
  role    = "roles/container.admin"
  member  = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}
