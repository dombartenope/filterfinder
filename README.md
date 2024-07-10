# Filter Finder
Script to allow the cross referencing of a user tag list string to each row of a csv containing a filters column

### Usage
1. Add `input.csv` containing an export of notifications (using OneSignal dashboard export from Sent Messages page, or using View Notifications/CSV Export API endpoints)
2. `go run main.go`
3. In the terminal, you'll be asked to provide the App Id, Api Key, and Onesignal ID to make a View User request to the OneSignal server (output saved in user.json).
4. The program will run against tags and language from that user and output any matches found in the terminal
`out.csv` should be generated for rows containing matches and the conditions that the user tags met, will be output to the terminal

### TODO
- Since this was built out of necessity, it is not equipped to handle complex filter relations (>/< unix timestamps) This will be added in the future
- Still need to add in the other fields in order to sort by more than just language and data tags
- Optimize the structure of this project so it can more easily be navigated and read

### Changes
~~If the filters look for language, you'll need to add the language as a data tag for the time being. This will be resolved when user exports are used instead of passing in a string~~
- This has now been updated. The new flow now includes the auto-generation of a .env file to store the app id, api key, and onesignal id used to make a request to the [View User]("https://documentation.onesignal.com/reference/view-user") endpoint. This will unmarshall the api response into a struct which can then be used to retrieve the different user properties that might be used in a notification's filters. 
