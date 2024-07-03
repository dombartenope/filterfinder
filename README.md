# Filter Finder
Script to allow the cross referencing of a user tag list string to each row of a csv containing a filters column

### Usage
1. Add `input.csv` containing an export of notifications (using OneSignal dashboard export from Sent Messages page, or using View Notifications/CSV Export API endpoints)
2. Grab the tags from a user on the dash and copy the full string with brackets to input into value for userTags
- ![image](https://github.com/dombartenope/filterfinder/assets/56173293/cf775b7e-1d7f-4774-a58a-4ec7f5dfb4c2)
3. `go run main.go`
`out.csv` should be generated for rows containing matches and the conditions that the user tags met, will be output to the terminal

### Considerations
Since this was built out of necessity, it is not equipped to handle complex filter relations (>/< unix timestamps) This will be added in the future
