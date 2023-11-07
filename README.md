# FFLogsFetcher
fflogsAPIから任意のデータをFetchするCLI

## Overview

## Requirement

## Usage

事前にfflogs apiからaccess_tokenを取得し、環境変数に登録してください。
```console
$ export FFLOGS_ACCESS_TOKEN=$(curl -u ${client_id}:${client_secret} -d grant_type=client_credentials https://ja.fflogs.com/oauth/token | jq -r '.access_token')
```
詳細は [Documentation](https://ja.fflogs.com/api/docs)を参考にしてください。

## Features
