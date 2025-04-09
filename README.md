# reacjira

reacjira = emoji reaction + jira ticket create

## Usage

### prepare

add a configuration like below to the `reacjira.toml` file.

```toml
[[reacjiras]]
emoji = "test_task"
project = "TEST"
issue_type = "Task"
epic_key = "FOO-1234"
description = """h2. foo
bar
"""
slack_message = ""
```

### creating jira tickets

if you add an emoji reaction (e.g. :test_task:) to slack messege, reacjira will create a jira ticket in "Test" Project as IssueType="Task".

![image](https://user-images.githubusercontent.com/1635885/111093000-0ff33200-857b-11eb-949b-357edda881d8.png)

![image](https://user-images.githubusercontent.com/1635885/111093007-14b7e600-857b-11eb-856f-2de92ee44c61.png)


## Development

```
cp config.toml config.secret.toml
$EDITOR config.secret.toml ## write your slack & jira token

make run
```


## Architecture

main -> bot -> service -> {config, slack, jira}

|package|role|
|---|---|
|main|An endpoint of this application.|
|bot|Controls communication with slack RTM|
|service|Provide a high level service by using slack, jira clients|
|slack|a client of slack|
|jira|a client of jira|
|config|It defines the config object|
