# check-websites
Check the status of your websites.

If we find any problems (the website is down or any status code not belonging to 200's family) we will send an email to the admin.

All settings are in .env.

Commands for Docker

Build Image
```
docker build -t check-websites .
```

Run App
```
docker run check-websites
```
