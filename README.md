# AWS SAM Demo for my Serverless Threat Modelling talk

## Previous work

This repo began life from here:

https://rtoch.com/posts/serverless-golang-with-lambda-and-dynamo/

https://github.com/CrazyRoka/todo-app-lambda

It has been refactored and made more modular for my purposes.  This is experimental code and not designed for production use.

> :warning: This code sample is provided AS-IS and use this code at your own risk.

## How to deploy to AWS

### Requirements

* [Go](https://go.dev/)
* The [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) for deployment
* [Make](https://www.gnu.org/software/make/)

## Deployment

1. Make sure you have an active AWS developer profile or role. When you deploy resources will be deployed to your account that will cost $$. Not much $$ but just be aware.
2. First time, run `make init` and answer the SAM deployment questions (defaults are fine)
3. For subsequent deployments a `make deploy` will be all you need. 

## Testing

Testing is manual (TODO: automated testing!).

The deployed API Gateway URL is output post deployment.

Create a new record:

```bash
curl -s -X POST https://saj7obvw34.execute-api.ap-southeast-2.amazonaws.com/Prod/user -H 'Content-Type: application/json' -d '{"name":"Mr Magoo","address":"1 Special Court", "passport":"PA12345"}' | jq
{
  "id": "5f0d79cc-6759-4ee8-9585-b5bd6d2efd3f",
  "name": "Mr Magoo",
  "address": "AQICAHis0XB5UoLmQQUKu10hsNPMMnuqV8GH0sUJs1T/tkX9YgHUJh1fYj7CSdxxyy09LHURAAAAijCBhwYJKoZIhvcNAQcGoHoweAIBADBzBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDG90LL7u44QpLkFkawIBEIBGaNuVDkYr+4+dEt8trZNxL3p0Dwk0l2H819FCVQ2fTYHQ2ntBGQWakfJhURjlQInmz89RdwRt4SRC08jZKCipC1eCXWaJiw",
  "status": false,
  "passport": "AQICAHjE4hKKqhO5lU44Hq+uIJee5elvxi+7sbsyqPd1afRcBQHvW2rapYsRGdQa4Zyi+3o9AAAAgTB/BgkqhkiG9w0BBwagcjBwAgEAMGsGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMv83Q6iS8DQWymh5UAgEQgD6O3xKh+LzGoKPQro9Px7v4U5M+CAzZuXn1DyAxVqVugYKQcE8GtU6/l5JuIi75W+QzDXn9wFLLb9xzgAZnjg"
}
```

Review a record (only address field is decrypted):

```bash
curl -s -X GET https://saj7obvw34.execute-api.ap-southeast-2.amazonaws.com/Prod/user/5f0d79cc-6759-4ee8-9585-b5bd6d2efd3f | jq
{
  "id": "5f0d79cc-6759-4ee8-9585-b5bd6d2efd3f",
  "name": "Mr Magoo",
  "address": "1 Special Court",
  "status": false,
  "passport": "AQICAHjE4hKKqhO5lU44Hq+uIJee5elvxi+7sbsyqPd1afRcBQHvW2rapYsRGdQa4Zyi+3o9AAAAgTB/BgkqhkiG9w0BBwagcjBwAgEAMGsGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMv83Q6iS8DQWymh5UAgEQgD6O3xKh+LzGoKPQro9Px7v4U5M+CAzZuXn1DyAxVqVugYKQcE8GtU6/l5JuIi75W+QzDXn9wFLLb9xzgAZnjg"
}
```

Review a record and decrypt passport field via the passport service:

```bash
curl -s -X GET https://saj7obvw34.execute-api.ap-southeast-2.amazonaws.com/Prod/passport/5f0d79cc-6759-4ee8-9585-b5bd6d2efd3f |jq
{
  "id": "5f0d79cc-6759-4ee8-9585-b5bd6d2efd3f",
  "name": "Mr Magoo",
  "address": "AQICAHis0XB5UoLmQQUKu10hsNPMMnuqV8GH0sUJs1T/tkX9YgHUJh1fYj7CSdxxyy09LHURAAAAijCBhwYJKoZIhvcNAQcGoHoweAIBADBzBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDG90LL7u44QpLkFkawIBEIBGaNuVDkYr+4+dEt8trZNxL3p0Dwk0l2H819FCVQ2fTYHQ2ntBGQWakfJhURjlQInmz89RdwRt4SRC08jZKCipC1eCXWaJiw",
  "status": false,
  "passport": "PA12345"
}
```







