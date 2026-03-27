# Sports Statistics Bot

A Telegram bot for tracking sports exercises. Supports Russian and English languages.

The bot allows you to:
- Add completed exercises via text or buttons
- View statistics for any period
- Quickly repeat frequent exercises with one tap (Quick-Add)

## Requirements

- Go 1.24+
- PostgreSQL
- Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))

## Installation and Setup

### 1. Clone

```bash
git clone https://github.com/DmiTryAgain/sports-statistics.git
cd sports-statistics
```

### 2. Configure

```bash
cp config/local.toml.dist config/local.toml
```

Open `config/local.toml` and fill in:

```toml
[Bot]
Token = "your-token-from-BotFather"
Name = "your_bot_name"
ReplyFormat = "markdown"
Debug = false
Timeout = "30s"

[Database]
Addr     = "localhost:5432"
User     = "postgres"
Database = "sport_statsrv"
Password = "postgres"
```

### 3. Create the database

```bash
createdb sport_statsrv
psql -d sport_statsrv -f docs/schema.sql
```

### 4. Build and run

```bash
make build
make run
```

The bot will start listening for Telegram updates. Logs are written to `app.log`.

## Usage

### Getting started

Send `/start` to the bot — a keyboard with three buttons will appear: **Add**, **Show**, **Help**.

All commands can be entered as text or via buttons.

### Adding exercises

Text format: `<command> <exercise> [parameters]`

Commands: `add`, `сделал`, `добавь`

**Examples:**

```
add pull-ups 15
add bench press 80kg 10
add jogging 5km 25min
add plank 90sec
add weight hold 40kg 30sec
add squats 60kg 12
add walking 3km
add jogging 30min
```

Parameters can be written together (`80kg`) or separately (`80 kg`). Supported units:
- Weight: `kg`, `g`, `lbs` (`кг`, `г`)
- Distance: `km`, `m` (`км`, `м`)
- Time: `h`, `min`, `sec` (`ч`, `мин`, `сек`)

Compound time values are summed: `1h 30min` = 90 minutes.

### Viewing statistics

Text format: `<command> <exercises> [period]`

Commands: `show`, `покажи`, `показать`

**Examples:**

```
show pull-ups for week
show all for today
show bench press squats for month
show pull-ups for year
show all
```

Periods: `today`, `yesterday`, `week`, `last week`, `week before last`, `month`, `last month`, `month before last`, `year`, `last year`, `year before last`. Without a period — all time.

Weekdays are also supported: `monday` (or `mon`), `tuesday` (or `tue`), `wednesday` (or `wed`), `thursday` (or `thu`), `friday` (or `fri`), `saturday` (or `sat`), `sunday` (or `sun`). If the weekday has already passed this week — shows that day this week; if not yet — shows that day last week; if today — shows from midnight to now.

Date formats are also supported: `15.03.2026` or a range `01.03.2026-15.03.2026`.

Example output:

```
exercise        weight  count     sets
pull-ups        -       22        2
bench press     60kg    10        1
bench press     80kg    18        2
```

### Help

```
help
help add
help show
```

### Button mode

Press **Add** — the bot will offer a list of exercises to choose from, then step by step request parameters (weight, count, distance, time) depending on the exercise type.

If you frequently do the same exercises, the bot will show **Quick-Add** buttons — add with one tap based on your history.

After adding via buttons, the bot will show a hint with a ready-made text command for quick copying and reuse.

## Supported exercises

**Bodyweight:** pull-ups, push-ups, dips, abs, squats, lunges, burpee, skipping rope, hyperextension, leg raise, muscle-ups

**Weighted:** bench press, deadlift, lat pulldown, leg press, preacher curl, shoulder press, bent-over row, dumbbell curl, leg extension, leg curl, seated row, chest fly, tricep pushdown, romanian deadlift, hip thrust, lateral raise, shrugs

**Cardio:** jogging, walking

**Timed:** plank, wall sit, hang, hollow hold, superman, side plank

**Timed + weighted:** weight hold

## Development

```bash
make test     # run tests
make fmt      # format code
make lint     # linter
```

Technical documentation: [docs/README.md](docs/README.md)
