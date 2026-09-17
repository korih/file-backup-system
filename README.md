# Homelab Backup

This service should require environment variables or secrets injected into it, which it uses to connect to AWS S3. With this, it will have a sqlite database that it uses to track state. It should also upload a WAL but i don't think we need this yet.

It will basically take all the files that it is defined to backup, and then back them up. It will first take the file, encrypt it based on a gpg key we define, then upload the file to S3. For a start we can have it work one at a time, where we only have once instance and it will only upload one file at a time. But hopefully I can make it more techinical so that we can have multiple instances and parallel uploads to make things faster.

I want the configuration to be defined in a configmap in the cluster, or in a known location, that the program can read to know how it should behave. 

Maybe part of same project create a simple WEBUI that i can access on tailnet that lets me upload a file to be encrypted and uploaded too. But an API is probably enough too

We should keep a hash of everyfile, if the checksum is different then there has been a change in the file and we can upload a new one.

somehow server and client should be synced, for start, we never delete a file unless we have a file that we have a different checksum so we delete to replace with new value

For now runs on a cronjob, in future make it a daemon like service that keeps it in sync

would be good to also define something to reverse the encryption

## TODO
- [ ] Start application
- [ ] Read files one by one on new run, app is stateless and runs on cronjob
- [ ] Keep track of what has been read
- [ ] Read files and calculate the checksum, compare with database
- [ ] If the file is new or changed, we upload it to S3 
- [ ] Encrypt file with GPG before we upload
- [ ] Should also compress the file to save space
- [ ] Upload the file to S3
- [ ] Update database, mark as done, and move on to next file
- [ ] Repeat until all files that have been configured are backed up, then die
