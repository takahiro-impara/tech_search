
#! /usr/bin/sh
unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
assume_role="arn:aws:iam::882275506731:role/Switch-from-599453524280-frit-whitelist-user-2"
res=$(aws sts assume-role --role-arn ${assume_role} --role-session-name AWSCLI-Session --profile frit)
echo $res | jq .

AWS_ACCESS_KEY_ID=$(echo $res | jq -r .Credentials.AccessKeyId)
AWS_SECRET_ACCESS_KEY=$(echo $res | jq -r .Credentials.SecretAccessKey)
AWS_SESSION_TOKEN=$(echo $res | jq -r .Credentials.SessionToken)

export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN