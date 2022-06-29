# usairnet-map

A map showing all US Air Net locations with links to their respective pages.

This can be used by skydivers or pilots to check aviation weather near them, without needing to know in advance which stations might be close.


Data is stored as a large JSON dictionary called data.json which looks like:

```
{
  "states": {
    ...
    "PA": {
      "code": "PA",
      "name": "Pennsylvania",
      "stations": {
        ...
        "KFML": {
          "code": "KFML"
          "name": "Franklin"
          "latitude": 41.38,
          "longitude": -79.87
        },
        ...
      }
    },
    ...
  }
}
```

## GoLang Setup
```
$ go mod init scrape.go
$ go mod tidy
$ go run scrape.go
```


## Deployment (from scratch)

1. Configure profile (matching `~/.aws/credentials`) and tfstate S3 bucket in main.tf
1. Run the terraform script to create bucket and user
1. Go to AWS console and create access keys for user
1. Save these into github repo secrets as `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`
1. Configure CNAME record to redirect `@` to your bucket endpoint (look up in AWS console)
1. Push changes or run github workflow, then go to http://usairnetmap.com and see if it all works!
