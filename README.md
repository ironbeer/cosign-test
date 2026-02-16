# cosign-test

## 検証手順
```shell
BUNDLE_PATH=main.go.bundle
ARTIFACT_PATH=main.go
CERTIFICATE_IDENTITY=https://github.com/ironbeer/cosign-test/.github/workflows/cosign.yml@refs/heads/main
CERTIFICATE_ISSUER=https://token.actions.githubusercontent.com

go run main.go $BUNDLE_PATH $ARTIFACT_PATH $CERTIFICATE_IDENTITY $CERTIFICATE_ISSUER

> Verified OK
```
