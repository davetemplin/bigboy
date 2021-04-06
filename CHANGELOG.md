# Changelog

## v1.3.1

- GO modules set up
- Continuous deployment/release workflow
- Unit test main with DB mocks
- Separate file for CLI flags parser
- Bugfix: nesting with MySQL by numeric value
- Bugfix: split by date for MySQL using layout for time parsing
- Bugfix: create directory if target doesn't include `/` in the end when using split
- Documentation coverage improved
- Examples added

## v1.3.0

- Changed progress "rows written" to update the previous row for more compact output
- Added unit tests for config and args
- Added `-c` flag to modify base config name/path
- Changed base config name from `config.json` to `bigboy.json`
- GitHub Actions added
- Bash build script added

## v1.2.0

- Bugfix slice of bytes for MySQL values

## v1.1.0

- First stable release
