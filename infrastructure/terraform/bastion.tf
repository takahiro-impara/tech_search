resource "aws_eip" "bastion" {
  instance = aws_instance.bastion.id
  vpc      = "true"
}

resource "aws_key_pair" "bastion" {
  key_name   = "${var.service_name}-bastion-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDVOwx0PR4nKDxHuclZi7FPRUZGhLZtidYUc2Yd/Y8Xg4VU41BwkR2qeN3QxsXyCdHwkaiCYsPqIYgc2KeWahxWuh4WCN79S5XyrfHczEUuEvfeTOPRmHuJ9IMwtMThJZeoB8j9zjHX7wdnUa92g1RKqzecBjMfudjTRjpqbvxHYfub9kzVKJR0Apbvjjx6/rogU4ZCCaOy3O4ay5J2N1hoxOQAFYzT4OpQznzb0Ix1PNfmP+sFISH5ET/QHI66AeEbfotHwS6A15fsMNFT9qbtQV1Qm1mKbR/HkC6l5rn6KyR9YBwJI2vjhk7mGHNVmX+qpBhoEXxV8SEHvAvDctyCnFp1EPBArSbxqyhuUabRI06M7UdP6TL7IsNQKLboujMjjgULJ7KXpJ6Ar2xFqGEDxhVtaaRGG4T8LKTowpXxDn1UqhnCCVzmghyVhuoQo5KH+lIiv/kabpz36TdC7nlVVlwh6lo37gZdxe3CuDDcvjVkV4tx2q68z+VpOtd/eVc="
}

data "template_file" "bastion" {
  template = file("userdata-template/bastion.tpl")

  vars = {
    MAINTENANCE_USER = var.maintenance_user
    INSTANCE_NAME    = "${var.region}-migration-bastion"
  }
}

resource "aws_instance" "bastion" {
  ami                         = "ami-0beaa649c482330f7"
  ebs_optimized               = false
  instance_type               = "t3.small"
  monitoring                  = false
  key_name                    = aws_key_pair.bastion.key_name
  subnet_id                   = data.aws_subnets.public.ids[0]
  associate_public_ip_address = true
  source_dest_check           = true
  user_data                   = data.template_file.bastion.rendered
  vpc_security_group_ids = [
    aws_security_group.bastion.id
  ]
  root_block_device {
    volume_type           = "gp2"
    volume_size           = 10
    delete_on_termination = false
  }
  tags = {
    Name = "${var.service_name}-bastion"
  }
}

resource "aws_security_group" "bastion" {
  name        = "${var.service_name}-bastion-sg"
  description = "${var.service_name}-Bastion operation Security Group"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group_rule" "bastion-rule-tcp22" {
  description = "SSH from ${var.service_name} vpc"
  type        = "ingress"
  from_port   = 22
  to_port     = 22
  protocol    = "tcp"

  cidr_blocks = ["126.227.57.97/32"]

  security_group_id = aws_security_group.bastion.id
}