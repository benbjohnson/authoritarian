Authoritarian
=============

## Overview

Command line utility for authorizing Twitter users to an application. This can be useful if you own a bunch of accounts (such as feed accounts) that need to be authorized to a Twitter application you run. For example, the [Scuttlebutt](https://github.com/benbjohnson/scuttlebutt) project used this to authorize all its feeds.


## Usage

To use authoritarian, you need to have you Twitter application's key and secret. From there, simply install and run the CLI:

```sh
$ go get github.com/benbjohnson/authoritarian
$ authoritarian -key $TWITTER_APP_KEY -secret $TWITTER_APP_SECRET
Listening on http://localhost:10000
```

Now you simply need to open your web browser to [http://localhost:10000](http://localhost:10000) and log in with the Twitter account you want to authorize. After you authorize the account, you'll be presented with a user auth key and secret. Copy these values to wherever you need them since they are not saved anywhere else.

*NOTE: PLEASE DO NOT PUT YOUR KEY & SECRET IN YOUR APPLICATION AND CHECK THEM IN!*

## Acknowledgements

This project is largely based on @kurrik's [examples in twittergo](https://github.com/kurrik/twittergo-examples). It's just wrapped up into a little CLI so it's easy to use.
