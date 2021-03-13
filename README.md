# reacjira

reacjira = emoji reaction + jira ticket create

## Usage

### prepare

add a configuration like `{emoji = "test_task", project = "TEST", issue_type = "タスク", description="created from reacjira."}` to `reacjira.toml`.

### creating jira tickets

if you add an emoji reaction (e.g. :test_task:) to slack messege, reacjira will create a jira ticket in "Test" Project as IssueType="タスク".

## Development

```
cp config.toml config.secret.toml
$EDITOR config.secret.toml ## write your slack & jira token

make run
```
