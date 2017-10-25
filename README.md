# CSV analysis

Find all CSV files in the same directory and run an aggregation over a column / dimension. Other dimension can filter results into the agg. Every CSV file is treated individully. The program is memory intensive, but can be managed by setting -t flag accordingly to the hardware and CSV files size.

Data is dumped to ElasticSearch
```
docker-compose -f db.yml up
```

```
go run main.go -c "Error Code" -p /Users/*/Downloads -t 4 -dimension_key Device -dimension_val Chromecast

```

```
 -c               dimension to run the aggregation on (default "Error Code")
 -output          dimension value match that will output data rows. Will run over the aggregation field. Will be outputed to csv file
 -p               path to directory to find CSVs on (default to same directory where binary is executed on)
 -t               number of gourotines will be spawn (default 8)
 -dimension_key   dimension to use as filter (default "Device")
 -dimension_val   dimension value (default "Chromecast")
```
Binary is available:
```
./csv_analysis  -c "Error Code"  -p /Users/*/Downloads -t 4 -dimension_key Device -dimension_val Chromecast
```

Example queries: 

```
./csv_analysis  -c "Error Code"  -p /Users/*/Downloads -t 4 -dimension_key Device -dimension_val Chromecast
```

Will show an aggregation over errors, including only those rows where Device=Chromecast


```
./csv_analysis  -c "Error Code" -output "3002" -p /Users/*/Downloads -t 4 -dimension_key "Device" -dimension_val "PC( MAC )"
```

Will show an aggregation over errors, including only those rows where Device=PC( MAC ). All rows that match previous conditions and Error Code=3002 will be outputed to a CSV file.

# Warning

Unexpected results can be found if some value in rows contains commas separated strings, for example if a columm stores JSON.
To avoid this bad behavior delete such columms.