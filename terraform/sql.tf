resource "google_sql_database_instance" "master" {
  name = "master-instance-2"
  region = "europe-west1"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "users" {
  name      = "users-db"
  instance  = "${google_sql_database_instance.master.name}"
}

resource "google_sql_user" "users" {
  name     = "admin"
  instance = "${google_sql_database_instance.master.name}"
  host     = "%"
  password = "${var.google_cloud_sql_admin_user_password}"
}
