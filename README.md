# `ttimer` - Terminal Timer

`ttimer` is a simple timer that counts down time left in a terminal window. If run on Mac, Windows, or desktop Linux, it will send silent system notifications at 90% and 100% completion.

## Installing

To install to your system you can use 

```
go get github.com/drgrib/ttimer
```

To make it accessible on the commandline as `ttimer`, assuming you've added `$GOPATH/bin` to your `$PATH`, you can use

```
cd $GOPATH/src/github.com/drgrib/ttimer
go build
go install
```

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

Or if you want a very specific duration, you can specify it using

```
ttimer -t 3h2m5s
```

Or if you want a very short time

```
ttimer -t 30s
```

## End Time Timing

Let's say you need to leave for the bus by 8:12 am, which is coming up in the next hour. You could simply enter

```
ttimer -t 812
```

and `ttimer` will automatically infer the next occurence of `8:12`, which is `am`. E.g.

```
== 812a Timer ==
23m29s
```

If you want to force it to set a timer for 8:12 pm, you could use

```
ttimer -t 812p
```

Resulting in something like

```
== 812p Timer ==
12h22m25s
```

If you want a timer for 3:00 pm, you could simply enter

```
ttimer -t 3p
```

All end time timers are set to time to align to zero seconds on the minute so they will change over to new minutes with the system clock.

## Parsing Rules

* All integers less than `100` will be interpretted as minutes
* Any strings fitting a call to [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration) will be interpretted as that duration. E.g. `1m30s` or `2h`.
* Any strings ending in `a`, `p`, `am`, or `pm` will be interpretted as times. E.g. `1p` or `930a`.
* Any integers greater than `100` will be interpretted as times. E.g. `242` will be interpretted as the next occurence of `2:42` and set to `am` or `pm`, whichever is soonest.

## Exiting

To exit the timer at any time, simply enter `q`.

## Timezone Setting

The default timezone for end times is `America/Los_Angeles`. This can be changed using the `-z` option

```
ttimer -t 9a -z UTC
```

Once set, `ttimer` will save the last `-z` option and use it as the default timezone in subsequent runs.

For reference, the list of official `tz` database timezones can be found [here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).
