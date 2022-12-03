#!/bin/bash -v
yum update -y

amazon-linux-extras enable redis6
yum clean metadata
yum install redis -y