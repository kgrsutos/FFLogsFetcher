# FFLogsFetcher
fflogsAPIから任意のデータをFetchするCLI

## Overview

現状は以下の機能のみ（後々追加予定）
- `query.yaml`に書かれたReportIDからDeath Eventを取得し、キャラクター名とDeathした時に受けたアビリティ名をファイル出力する

## Requirement
- Go 1.20

## Usage

事前にfflogs apiからaccess_tokenを取得し、環境変数に登録してください。
```console
$ export FFLOGS_ACCESS_TOKEN=$(curl -u ${client_id}:${client_secret} -d grant_type=client_credentials https://ja.fflogs.com/oauth/token | jq -r '.access_token')
```
詳細は [Documentation](https://ja.fflogs.com/api/docs)を参考にしてください。

以下コマンドを実行すると、output.txtが出力されます。
```console
$ go run main.go
```

## Reference
- https://ja.fflogs.com/v2-api-docs/ff/
