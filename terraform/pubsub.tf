resource "google_pubsub_topic" "eventd-ingestion" {
  name = "event-ingestion"
}

resource "google_pubsub_subscription" "required-persistant-subscription" {
  name  = "required-persistant"
  topic = "${google_pubsub_topic.eventd-ingestion.name}"

  ack_deadline_seconds = 30
}
