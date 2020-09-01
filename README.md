# Nalupi
Compute PI to infinite digits, similar to y-cruncher

## Basic Usage
```bash
go run cmd/nalupi/main.go --decimal=100
```
If you have go installed, you can simply get PI value to any precision with the above command.

Alternatively, you can get the last-computed PI precision from server from [here](https://nalupi-b235sdkoha-de.a.run.app/pi/current)

The estimation of Sun circumference using the pre-computed PI can also be retrieved [here](https://nalupi-b235sdkoha-de.a.run.app/sun/circumference)

Note that the server runs every 5 minute through an app-script timer to compute the next precision

## Spreadsheet references
To prevent server from crashing or any potential issue saving snapshot and breaking the correctness of the computation (rate limit by Google Spreadsheet), we suggest you to view the values of PI or snapshots directly through the spreadsheet links down below
- [Last computed PI](https://docs.google.com/spreadsheets/d/1YnXZwX5ABPmBUFhktGVLDVnmgluVgSMFjIkMyIJ8Lt0/edit#gid=0)
- [Last snapshot](https://docs.google.com/spreadsheets/d/1FMUFV2z_MaccKswNLh3-x2vDeBY3RRNNzzAusjh848c/edit#gid=0)
- [Chudnovsky fraction metadata snapshot](https://docs.google.com/spreadsheets/d/1w7yT7uS-JmvvF9flQRQjqiX18bd9c0I30B-4x7EHLVw/edit)

## References
- [Pi in the Sky Youtube](https://youtu.be/BwkpNd2ceBk)