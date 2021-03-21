# reacjira

reacjira = emoji reaction + jira ticket create

## Usage

### prepare

add a configuration like `{emoji = "test_task", project = "TEST", issue_type = "タスク", description="created from reacjira."}` to `reacjira.toml`.

### creating jira tickets

if you add an emoji reaction (e.g. :test_task:) to slack messege, reacjira will create a jira ticket in "Test" Project as IssueType="タスク".

![image](https://user-images.githubusercontent.com/1635885/111093000-0ff33200-857b-11eb-949b-357edda881d8.png)

![image](https://user-images.githubusercontent.com/1635885/111093007-14b7e600-857b-11eb-856f-2de92ee44c61.png)


## Development

```
cp config.toml config.secret.toml
$EDITOR config.secret.toml ## write your slack & jira token

make run
```
