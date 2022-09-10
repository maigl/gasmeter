# Gasmeter

I want to track my gas consumption. Interesting data could be absolute
consumption and also daily consumption. Maybe also a current consumption could
be useful. Eventually, these values shall be displayed in some sort of
dashboard.

My gasmeter emits magnetic pulses for each 0.01 m3 gas. I want to use a
raspberry and a simple app to track those pulses.

## Goals

* Each magnetic pulse is detected by the app and stored in a database the
* current absolute value can not be read from the gasmeter, we can only count
  increments
* the app must allow manual updates of the current absolute gasmeter value

## Design

Each pulse is stored in a postgres database (postres already exists) with timestamp in unix format and in human readable format and
the last known gasmeter value incremented by the increment.

The gasmeter value is counted as int to avoid floating point issues, i.e the unit is 0.01 m3. The value will be
converted to m3 in all r/w operations.

There is http server with a GET endpoint that returns the current gasmeter value and a POST endpoint that allows
to perform a manual update.

## Issues

Most probably, the detection of the pulses will not be perfect and so we will deviate from the
real gasmeter value. Testing is a little tricky, the idea is to make manual tests and see
how bad it is.

## Deployment

Currently the app is deployed on a raspberry pi 3. The app is started with systemd.

```bash
systemctl --user status gasmeter
```
