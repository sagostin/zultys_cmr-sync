### Compilation

1. Download the Sourcecode
2. Install the latest version of Go
3. From CLI, run `go build`
4. Open the compiled program using your CLI of choice.

### Running

1. After compiling the program, run the program from CLI, using the flags.
   eg. ./hubspot-call_contact --apikey=YOUR_API_KEY --recent
   Including the `--recent` flag will only pull recent calls, if no calls have been recently processed, it will grab all calls
2. You may ignore the `--recent `flag to process all calls, every time added

### Installing (Linux)

1. Copy the binary file & storage.json or make a storage.json

   ```
   {
    "latest_call_timestamp": "2022-11-09T23:53:16.805Z"
   }
   ```
2. Add the path to the binary w/ flags to your cron file.

   `*/5 * * * * bash /opt/call_contact/call-contact_linux --apikey=APIKEY --recent`
