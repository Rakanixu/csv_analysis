#CSV analysis

Find all CSV files in the same directory and run an aggregation over a column / dimension. Other dimension can filter results into the agg.
Every CSV file is treated individully. 

go run main.go -c "Error Code" -dimension_key Device -dimension_val Chromecast

-c               dimension to run the aggregation on (default "Error Code")
-dimension_key   dimension to use as filter (default "Device")
-dimension_val   dimension value (default "Chromecast")

Binary is available:
./csv_analysis -c "Error Code" -dimension_key Device -dimension_val Chromecast
