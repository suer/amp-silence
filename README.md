# amp-silence

This is a simple CLI to add/delete silences in AMP(Amazon Managed Service for Prometheus) Alertmanager.

## Usage

### Add silence

```bash
$ echo '{
    "startsAt": "2024-06-30T15:00:00.000Z",
    "endsAt": "2024-06-30T23:59:59.000Z",
    "comment": "Maintenance",
    "createdBy": "suer",
    "matchers": [        {
            "name": "host",
            "value": "www.example.com",
            "isEqual": true,
            "isRegex": false
        }
    ]
}' | ./amp-silence add -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/

{"silenceID": "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"}
```

Required permission: `aps:PutAlertManagerSilences`

### Delete silence

```bash
$ ./amp-silence delete -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/ -s yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy
```

Required permission: `aps:DeleteAlertManagerSilence`

### List silences

```bash
$ ./amp-silence list -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/
[{"comment":"Maintenance","createdBy":"suer","endsAt":"2024-06-30T23:59:59.000Z","id":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","matchers":[{"isEqual":true,"isRegex":false,"name":"host","value":"www.example.com"}],"startsAt":"2024-06-30T15:00:00.000Z","status":{"state":"pending"},"updatedAt":"2024-06-29T10:57:20.518Z"}]
```

Required permission: `aps:ListAlertManagerSilences`

### List workspaces

```bash
$ ./amp-silence list-workspaces
[{"Arn":"arn:aws:aps:ap-northeast-1:XXXXXXXXXXXX:workspace/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","CreatedAt":"2024-02-06T08:33:51.749Z","Status":{"StatusCode":"ACTIVE"},"WorkspaceId":"ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","Alias":"aaaa","KmsKeyArn":null,"Tags":{}}]
```

Required permission: `aps:ListWorkspaces`

## Advanced Usage

### Query with JMESPath

```bash
$ ./amp-silence list -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/ -q "[0].id"
"yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
```

### Raw output

```bash
$ ./amp-silence list -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/ -q "[0].id" -r 
yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy
```

### Request timeout

The HTTP request timeout is 30 seconds by default. Set `AMP_SILENCE_TIMEOUT` (Go duration format) to change it.

```bash
$ AMP_SILENCE_TIMEOUT=10s ./amp-silence list -e https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/
```

## Build

```bash
$ go build
```
