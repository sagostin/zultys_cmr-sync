### Compilation

1. Download the Sourcecode
2. Install the latest version of Go
3. From CLI, run `go build`
4. Open the compiled program using your CLI of choice.

### Running

1. After compiling the program, run the program from CLI, using the flags.
   eg. ./hubspot-call_contact --apikey=YOUR_API_KEY --recent
   Including the `--recent` flag will only pull recent calls, if no calls have been recently processed, it will grab all calls
2. You may ignore the `--recent `flag to process all calls, every time.added
