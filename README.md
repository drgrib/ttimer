# Terminal Timer

`ttimer` is a simple timer that counts down time left in a terminal window. It will send silent system notifications (on mac, windows, or desktop linux) at 90% and 100% completion of the timer.

## Duration Timing

Lets say you want a timer for 3 minutes. Simply enter

```
ttimer -t 3
```

and you will start a timer count down like so

```
== 3m Timer ==
2m55s
```

## Parsing Rules

* All integers less than `100` will be interpretted as minutes
* Any strings fitting a call to [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration) will be interpretted as that duration. E.g. `1m30s` or `2h`.
* Any strings ending in `a`, `p`, `am`, or `pm` will be interpretted as times. E.g. `1p` or `930a`.
* Any integers greater than `100` will be interpretted as times. E.g. `242` will be interpretted as the next occurence of `2:42` and set to `am` or `pm`, whichever is soonest.

## Timezone Setting

The default timezone for time-based timers is `America/Los_Angeles`. It can be changed using the `-z` option

```
ttimer -t 9a -z UTC
```

Once set, `ttimer` will save the last `-z` option and use it as the default timezone in subsequent runs.

For reference, the list of official `tz` database timezones can be found [here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).
