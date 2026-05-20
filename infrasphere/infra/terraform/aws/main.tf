variable "region" { default = "us-west-2" }
provider "aws" { region = var.region }
output "note" { value = "Starter module placeholder for ECS/EKS/RDS-backed InfraSphere deployment." }

