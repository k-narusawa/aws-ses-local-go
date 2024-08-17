# aws-ses-local-go

![coverage](https://raw.githubusercontent.com/k-narusawa/aws-ses-local-go/badges/.badges/main/coverage.svg)

## cli

### v1

```shell
aws ses send-email \
--to from@example.com \
--from to@example.com \
--subject "subject" \
--text "body" \
--endpoint-url "http://localhost:8080"
```

```shell
aws ses send-raw-email \
--raw-message file://./scripts/message.json \
--endpoint-url "http://localhost:8080"
```
