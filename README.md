![onehit logo](https://github-production-user-asset-6210df.s3.amazonaws.com/96031819/273399434-5f5066f4-507e-4333-9378-3acb765ef5ab.png)
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# onehit

**onehit** is a hit counter as a service. It provides an HTTP API that allows you to increment and track hits for any arbitrary "key". This can be useful for tracking page views, downloads, or any other metric you want to count.

## how it Works

**onehit** is built on top of [libkeva](https://github.com/donuts-are-good/libkeva), a lightweight key-value store library for Go. libkeva provides a thread-safe API for storing and retrieving arbitrary data, which makes it a perfect fit for a hit counter service like **onehit**.

When a `GET` request is made to a path starting with `/x/`, **onehit** increments the hit count for the key specified in the rest of the path. The current hit count for the key is then returned in the response.

## usage

To use **onehit**, simply start the service and make a GET request to the `/x/{your-key}` endpoint, replacing `{your-key}` with the key you want to track. For example:

```
curl http://localhost:3589/x/my-page
```

This will increment the hit count for "my-page" and return the current count.

## deployment

For easy deployment with automatic HTTPS support, try serving **onehit** with [appserve](https://github.com/donuts-are-good/appserve), an easy application server that handles HTTPS for you automatically.


## license

mit license 2023 donuts-are-good