module "networking" {
  source = "./modules/networking"
  region               = "us-east-1"
  environment          = "prod"
  vpc_cidr             = "10.0.0.0/16"
  ec2-ami = "ami-0ff8a91507f77f867"
  max-subnets = 1
}

output "ssh" {
  description = "Connection data for new instance"
  value       = module.networking.ssh
}