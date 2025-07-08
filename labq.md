# LabQ Additions

This project now supports defining administrators and listing subjects via Telegram commands.

## Configuration

Add a new `admins` field to your config YAML with Telegram user IDs allowed to manage classes and schedules:

```yaml
admins:
  - 123456789
```

## Commands

- `/subjects` - list all available subjects with their IDs
- `/add_class <name>` - create a new class (admin only)
- `/add_date <subject_id> <day_of_week> <time(HH:MM)> <interval_weeks> <start_date(YYYY-MM-DD)>` - add schedule entry (admin only)
- `/join <subject_id>` - join the queue for the subject
- `/queue <subject_id>` - show the queue
- `/setname` - set or change your displayed name

Administrators are recognised based on the IDs configured above. Regular users can join and view queues.
