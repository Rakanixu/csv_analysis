# CSV analysis

Find all CSV files in the same directory and run an aggregation over a column / dimension. Other dimension can filter results into the agg. Every CSV file is treated individully. The program is memory intensive, but can be managed by setting -t flag accordingly to the hardware.

```
go run main.go -c "Error Code" -p /Users/*/Downloads -t 8 -dimension_key Device -dimension_val Chromecast

```

```
 -c               dimension to run the aggregation on (default "Error Code")
 -p               path to directory to find CSVs on (default to same directory where binary is executed on)
 -t               number of gourotines will be spawn (default 8)
 -dimension_key   dimension to use as filter (default "Device")
 -dimension_val   dimension value (default "Chromecast")
```
Binary is available:
```
./csv_analysis  -c "Error Code" -p /Users/*/Downloads -t 8 -dimension_key Device -dimension_val Chromecast
```
