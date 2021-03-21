terraform {
  backend "s3" {
    profile = "arbeiter"
    bucket = "intense-tfstates"
    key = "usairnetmap"
    region = "us-east-2"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
}


provider "aws" {
  profile = "arbeiter"
  region  = "us-west-2"
}


resource "aws_s3_bucket" "usairnetmap" {
  bucket = "usairnetmap.com"
  acl = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}


resource "aws_s3_bucket_policy" "usairnetmap" {
  bucket = aws_s3_bucket.usairnetmap.id
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "PublicReadGetObject",
        "Effect": "Allow",
        "Principal": "*",
        "Action": [
          "s3:GetObject"
        ],
        "Resource": [
          "${aws_s3_bucket.usairnetmap.arn}/*",
        ]
      }
    ]
  })
}


resource "aws_iam_user" "usairnetmap" {
  name = "usairnetmap"
}


resource "aws_iam_user_policy" "usairnetmap" {
  name = "usairnetmap"
  user = aws_iam_user.usairnetmap.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        "Effect": "Allow",
        "Action": [
          "s3:PutObject",
        ],
        "Resource": "arn:aws:s3:::usairnetmap.com/*",
      },
    ]
  })
}


# vim: ts=2 expandtab
