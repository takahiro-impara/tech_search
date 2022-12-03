data "aws_subnets" "private" {
  filter {
    name   = "tag:Name"
    values = ["education-vpc-private*"] # insert values here
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "tag:Name"
    values = ["education-vpc-pub*"] # insert values here
  }
}