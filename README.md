# Should I water the trees?

Checks recent precipiation data in your area and decides whether it's time to water unestablished trees which are less than two years old.

[Guide for watering trees](https://arbordayblog.org/treecare/how-to-properly-water-your-trees/)

## Setup

First, [obtain a token](https://www.ncdc.noaa.gov/cdo-web/token)

Next, copy the `config-template.yaml` into a new file `config.yaml` and update the file to contain your token.

## Running

You need to pass the station ID as an argument into the program. [Find your station ID here](https://www1.ncdc.noaa.gov/pub/data/ghcn/daily/ghcnd-stations.txt)

For example, the station ID for Central Park is `USW00094728`

Example usage:
```
go run shouldwater.go USW00094728
```

## Known Issues

The NCDC data is not always up to date and sometimes days are missing. However, when you retrieve data through another way it can contain data that NCDC doesn't have. Example request: 

```
https://www.ncei.noaa.gov/access/services/data/v1?dataset=daily-summaries&startDate=2022-01-29&endDate=2022-02-04&stations=USW00094728&format=json
```
