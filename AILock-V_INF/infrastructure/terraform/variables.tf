variable "project_id" {
  description = "The project ID for GCP"
  type        = string
  default     = "ailock-project"
}

variable "region" {
  description = "The region for GCP resources"
  type        = string
  default     = "us-central1"
}
