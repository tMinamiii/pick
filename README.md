# Pick

 Pick is able to search sites in your Pocket and open browser.

[![Image from Gyazo](https://i.gyazo.com/5ed83885a636c52ba43ccb2527002f90.gif)](https://gyazo.com/5ed83885a636c52ba43ccb2527002f90)

## Installation

``` sh
go get -u github.com/tMinamiii/pick
```

## Usage

### Create Pocket Cosumer Key

Access https://getpocket.com/developer/apps/ and **CREATE AN APPLICATION**

1. input `Application Name`(e.g. pick)
2. Check Permission `Retrieve`
3. Check Platform `Desktop (other)`
4. Check `I accept the Terms of Service`
5. Push **CREATE APPLICATION**, then generate consumer key.

### Auth Pocket

``` sh
pick auth <Pocket Consumer Key>
```

Generate authorization token file in `$HOME/.config/pick/key.json`.

### Run Pick

``` sh
pick
```

## Motivation

 Pocket is one of very useful service, however we keep storing favorite
or read later site.  We often only add site and never visit it. Pick helps
search titles ( and contets if you are premium user ) in your Pocket sites and
open browser directly.
