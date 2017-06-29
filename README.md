#CSV analysis

Find all CSV files in the same directory and run an aggregation over a column / dimension. Other dimension can filter results into the agg.
Every CSV file is treated individully. 

go run main.go -c "Error Code" -device_key Device -device_val Chromecast

-c            dimension to run the aggregation on (default "Error Code")
-device_key   dimension to use as filter (default "Device")
-device_val   dimension value (default "Chromecast")

Binary is available:
./csv_analysis -c "Error Code" -device_key Device -device_val Chromecast
