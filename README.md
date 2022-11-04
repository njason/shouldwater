# Should I water the trees? 🌳🌧️

A low cost solution that uses forecast and recent historical weather data to alert tree stewards when watering is needed for unestablished trees which are less than two years old. This solutions uses tomorrow.io for weather data and mailchimp for emailing alerts.

[Guide for watering unestablished trees](https://vimeo.com/416031708#t=5m35s)

## Setup


### Config

Copy the `config-template.yaml` into a new file `config.yaml`, and update the following fields:

- `tomorrowioApiKey`: After creating a free [tomorrow.io](https://www.tomorrow.io/) account, find the `Secret Key` [here](https://app.tomorrow.io/development/keys).
- `mailchimp.apiKey`: After creating a free [mailchimp](https://mailchimp.com/) account, create an API key [here](https://admin.mailchimp.com/account/api/)
- `mailchimp.templateId`: [Create a template](https://mailchimp.com/help/create-a-template-with-the-template-builder/) to use for alerting tree stewards to water. [This](https://us13.admin.mailchimp.com/templates/share?id=174361973_a7f368481da096f6c0df_us13) is the template used in NYC you can use as a starting point.
- `mailchimp.listId`: Use [this guide](https://mailchimp.com/help/find-audience-id/) to find the list/audience ID.
- `lat`, `lng`: The coordinates of where to run. You can use [Google Maps](https://support.google.com/maps/answer/18539) to find coordinates, format `lat, lng`.

### Build

To build a binary, run this from the repo root directory:

```
go build .
```

### Deploying

The binary will need to be in the same directory as `config.yaml` to run.

The binary is intented to run in two modes, the default action of determining if the alert email should go out, or if data should be recorded.

#### Record mode

Record every 6 hours, which is the limit of free historical data from Tomorrow.io. For macOS and linux, use [cron](https://phoenixnap.com/kb/set-up-cron-job-linux) with a line like `59 0,6,12,18 * * * /path/to/shouldwater record`. On Windows 10 use [Task Scheduler](https://www.windowscentral.com/how-create-automated-task-using-task-scheduler-windows-10).

#### Default mode

Once a week run in default mode to decide if an alert goes out. It is recommended to run default mode on Fridays so tree stewards can water on the weekends. The cron for this would be `0 16 * * FRI /path/to/shouldwater record`.
