#####################################################################
# Google Cloud Platform
#####################################################################
provider "google" {
    project = "${var.google_cloud_project_id}"
    region = "europe-west1"
}
