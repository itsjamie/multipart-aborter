# multipart-aborter
Just a small helper golang to abort all pending multipart uploads in a Amazon S3 bucket.

### Install
```
go get github.com/itsjamie/multipart-aborter
```

Example usage:
```
AWS_ACCESS_KEY_ID=? AWS_SECRET_ACCESS_KEY=? AWS_REGION=? multipart-aborter --bucket ?
```
