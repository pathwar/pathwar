#!/usr/bin/env bash
 
# This is a script shell that validates a Pathwar coupon
# It uses the pathwar cli tool, available at:
#
# - https://github.com/pathwar/node-pathwar
 
# this is the coupon you want to validate
coupon=$1
 
# this is *your* pathwar token, you can fetch it by analyzing the HTTP requests on the portal
# we will soon provide a page on your account to display the token.
export pathwar_token=
 
# fetch the id of the Epitech session
session=$(pathwar ls sessions name=Epitech2015 -q)
 
# get the id of your team in the Epitech session
organization=$(pathwar ls teams session=$session -q)
 
# validate the coupon for your team
pathwar add organization-coupons organization=$organization coupon=$coupon
