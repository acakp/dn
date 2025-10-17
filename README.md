# dn - Create Daily Note

Type `dn` to open the daily note at ~/ in your default editor (or in notepad.exe on windows).

## Configuration

The `dn` command searches for `/dn/config.toml` in the os's default config directory. On Linux, it will look for `~/.config/dn/config.toml`, and on Windows, for `%AppData%\dn\config.toml`.

In the config, the following settings can be defined:

### `path`

Path where the note will be opened. Example usage:

```
path = "~/notes/daily notes"
```

### `editor`

Command that will be passed to the command line to open the daily note file. Example usage:

```
editor = "vi"
```

### `extension`

File extension for the daily note. Example usage:

```
extension = "txt"
```

### `format`

How the note will be named. You can customize the naming by using the format strings in the following table:

| %     | Example Output |
|-------|----------------|
| %YYYY | 2025           |
| %YY   | 25             |
| %MM   | 10             |
| %M    | October        |
| %D    | 17             |
| %WW   | 05             |
| %W    | Friday         |
| %w    | Fri            |

## Flags

All settings from the config can be passed as flags. Type `dn -h` for more info.
Type `dn -1` to open a note from a previous day. This command opens a note named with the date of the day preceding the current one. For example, if the format is "%YYYY-%MM-%D" and today's date is "2025-10-17", using the `-1` flag will yield "2025-10-16". If the format is "%YYYY", using the `-1` flag will yield "2024".
