#CSV analysis

Find all CSV files in the same directory and run an aggregation over a column / dimension. Other dimension can filter results into the agg.
Every CSV file is treated individully. 

go run main.go -c "Error Code" -p /Users/*/Downloads -m 800000000 -dimension_key Device -dimension_val Chromecast

-c               dimension to run the aggregation on (default "Error Code")
-p               path to directory to find CSVs on (default to same directory where binary is executed on)
-m               bytes allocated to read CSV file. (default  800MB)
-dimension_key   dimension to use as filter (default "Device")
-dimension_val   dimension value (default "Chromecast")

Binary is available:
./csv_analysis  -c "Error Code" -p /Users/*/Downloads -m 800000000 -dimension_key Device -dimension_val Chromecast
