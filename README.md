# Should I water the trees? 🌳⛆

A low cost solution that uses forecast and recent historical weather data to decide whether it's time to water unestablished trees which are less than two years old.

[Guide for watering trees](https://arbordayblog.org/treecare/how-to-properly-water-your-trees/)

## Setup

Copy the `config-template.yaml` into a new file `config.yaml`, and update the following fields:

- `tomorrowioApiKey`: After creating a free [tomorrow.io](https://www.tomorrow.io/) account. Find the `Secret Key` [here](https://app.tomorrow.io/development/keys).
- `lat`, `lng`: The coordinates of where to run. You can use [Google Maps](https://support.google.com/maps/answer/18539) to find coordinates, format `lat, lng`.

Build a binary

linux, macOS:

```
go build -o shouldwater-cli .
```

Windows

```
go build .
```

The binary will need to be in the same directory as `config.yaml` to run.

### Deploying

The binary is intented to run in record mode every 6 hours, which is the limit of free historical data from Tomorrow.io. For macOS and linux, use [cron](https://phoenixnap.com/kb/set-up-cron-job-linux) with a line like `59 0,6,12,18 * * * /path/to/shouldwater record`. On Windows 10 use [Task Scheduler](https://www.windowscentral.com/how-create-automated-task-using-task-scheduler-windows-10).

After a week's worth of data has been collected (168 hourly records), run `shouldwater` to determine if watering is needed.
