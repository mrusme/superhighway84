Superhighway84
--------------

[![Superhighway84](superhighway84.jpeg)](superhighway84.png)

```
===============================================================================
                        INTERACTIVE ASYNC / FULL DUPLEX
===============================================================================

                              Dial Up To 19.2 Kbps
                                  
                                      with

    _  _ _ __ ____                  __   _      __                   ___  ____
   /  / / // / __/_ _____  ___ ____/ /  (_)__ _/ / _    _____ ___ __( _ )/ / /
  _\ _\_\_\\_\ \/ // / _ \/ -_) __/ _ \/ / _ \/ _ \ |/|/ / _ \/ // / _  /_  _/
 /  / / // /___/\_,_/ .__/\__/_/ /_//_/_/\_, /_//_/__,__/\_,_/\_, /\___/ /_/
                   /_/                  /___/                /___/

 ::: UNCENSORABLE USENET-INSPIRED DECENTRALIZED INTERNET DISCUSSION SYSTEM :::


   The V.H.S. (Very High Speed) Superhighway84 platform is more than just the
  fastest decentralized, uncensorable, USENET-inspired communications platform 
         available. It is also the first one to be based on the latest 
                        IPFS technology available today!

    Superhighway84 offers the most spectacular features under the Spectrum.
                                       
                             100% Error Protection
                         Data and Character Compression
                         Alternate Bell Compatible Mode
                         Long Haul Satellite Operation
                              Network Diagnostics
                                 Fallback Mode
                                   And More!


                    The Superhighway84 modern, uncensorable, 
                   decentralized internet discussion system.
                       It should cost a lot more than $0.


```

![Screenshot](screenshot01.png)

Superhighway84 is an open source, terminal-based, IPFS-powered, USENET-inspired,
uncensorable, decentralized peer-to-peer internet discussion system with retro
aesthetics.

[More info here.](https://xn--gckvb8fzb.com/superhighway84/)



## Installation

### Prerequisites

Download the [kubo 0.16
release](https://github.com/ipfs/kubo/releases/tag/v0.16.0) and unpack it:

```sh
$ tar -xzf ./kubo_*.tar.gz
```

If you haven't used IPFS so far, initialize the IPFS repository using the 
following command:

```sh
$ ./kubo/ipfs init
```

If you had used IPFS an already have an IPFS repository in place, either
(re)move it from `~/.ipfs` or make sure to `export IPFS_PATH` before running the
`ipfs init` command, e.g.:

```sh
$ export IPFS_PATH=~/.ipfs-sh84
$ ./go-ipfs/ipfs init
```


### From Release

Download the [latest
release](https://github.com/mrusme/superhighway84/releases/latest) and unpack
it:

```sh
$ tar -xzf ./superhighway84_*.tar.gz
$ ./superhighway84
```

If you initialized the IPFS repo under in a custom location, you need to prefix
`IPFS_PATH`:

```sh
$ IPFS_PATH=~/.ipfs-sh84 ./superhighway84
```

The binary `superhighway84` can be moved wherever you please.



### From Source

Clone this repository

- from [GitHub](https://github.com/mrusme/superhighway84)
  ```sh
  $ git@github.com:mrusme/superhighway84.git
  ```
- from [Radicle](https://app.radicle.network/seeds/maple.radicle.garden/rad:git:hnrkcf9617a8pxxtw8caaop9ioe8cj5u4c4co)
  ```sh
  $ rad clone rad://maple.radicle.garden/hnrkcf9617a8pxxtw8caaop9ioe8cj5u4c4co
  ```

Then cd into the cloned directory and run:

```sh
$ go build .
```

The binary will be available at ./superhighway84 and can be moved wherever you
please.



## Running

First, check ulimit -n and verify that it's at a reasonable amount. IPFS
requires it to be large enough (>= 2048) in order to work properly over time.

Second, if your hardware shouldn't be a beefy computer but instead one of
those flimsy MacBooks, older hardware, a Raspberry or a low-memory VPS it is
advisable to set the previously created IPFS repository to the `lowpower`
profile.

```sh
$ ipfs config profile apply lowpower
```

This should help with CPU usage, file descriptors and the amount of network
connections. While during the startup period you might still see peers peaking
between 1k and 3k, connections should ultimately settle somewhere between 100
and 300 peers.

Afterwards you can simply launch the binary:

```sh
$ superhighway84
```

A setup wizard will help you with initial configuration. Please make sure to
have at least HOME and EDITOR exported in your environment.

In case you're intending to run the official IPFS daemon and Superhighway84 in
parallel, be sure to adjust the ports in their respective IPFS repos (e.g.
`~/.ipfs` and `~/.ipfs-sh84`) so that they won't utilize the same port numbers.
The ports `4001`, `5001` and `8080` are relevant and should be adjusted to
something other for every new repo/IPFS node that will run in parallel, e.g.:

```json
  "Addresses": {
    "Swarm": [
      "/ip4/0.0.0.0/tcp/4002",
      "/ip6/::/tcp/4002",
      "/ip4/0.0.0.0/udp/4002/quic",
      "/ip6/::/udp/4002/quic"
    ],
    "Announce": [],
    "NoAnnounce": [],
    "API": "/ip4/127.0.0.1/tcp/5002",
    "Gateway": "/ip4/127.0.0.1/tcp/8081"
  },
```

**NOTE**: When running Superhighway84 for the first time it might seem like it's
"hanging" at the command prompt. Usually it isn't hanging but rather searching
for peer it can connect to in order to synchronize the database. Depending on
how many people are online, this process might take _some time_, please be
patient.



## Connectivity

If you're having trouble connecting to the IPFS network that might be due to
your network setup. Please try the IPFS `AutoRelay` feature in such a case:

```sh
$ ipfs config --json Swarm.RelayClient.Enabled true
```

More information on this can be found here:
https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#autorelay



## Configuration

Superhighway84 will guide you through the basic configuration on its first run.
The configuration is stored at the path that you specified in the setup wizard.
After it was successfully created, it can be adjusted manually and will take
effect on the next launch of Superhighway84.

Configuration options that might be of interest:

```
ArticlesListView =
  The view to be used for the articles lit. Possible values:
  0 - threaded view, latest thread at the top
  1 - list view, latest article at the top

[Profile]
  From =
    The identifier that is being shown when posting an article, e.g. your name,
    username or email that you'd like to display

  Organization =
    An optional organization that you'd like to display affiliation with

[Shortcuts]
  The shortcuts for navigating Superhighway84, can be reset to its defaults by
  simply removing the whole [Shortcuts] block and launching Superhighway84

  The structure is as following:

  `<key code> = "event"`

  The key codes can be looked up under the following link:

  https://pkg.go.dev/github.com/gdamore/tcell/v2#Key

  For simple ASCII characters use their ASCII code, e.g. `114` for the character 
  `r`.
```


## Usage

The default keyboard shortcuts are:

```
     C-r: Refresh
     C-h: Focus groups list
C-l, C-k: Focus articles list
     C-j: Focus preview pane
     C-q: Quit
       k: Move up in list
       j: Move down in list
       h: Move left in list
       l: Move right in list
       g: Move to the beginning of list/text
       G: Move to the end of list/text
      CR: Select item in list
       n: Publish new article
       r: Reply to selected article
```

However, you are free to customize these within your configuration file, under
the section `Shortcuts`. 


### Submit Article

When submitting a new article or a reply to an article, the $EDITOR is launched
in which a document with a specific structure will be visible. This structure
consists of the HEADER, a SEPARATOR and the BODY and looks like this:

```
Subject: This is the subject of the article
Newsgroup: test.sandbox
= = = = = =
This is the multiline
body of the article
```

The HEADER contains all headers that are required for an article to be
submitted. These are:

- `Subject:`\
  The subject of the article that will be shown in the articles list. The
  subject must only contain of printable ASCII characters.

- `Newsgroup:`\
  The newsgroup under which the article will be submitted, this can
  either be an existing group or a new group. Please try to follow
  the convention when creating new groups.
  The newsgroup must only contain of printable ASCII characters.

The SEPARATOR contains of 6 equal signs and 5 spaces, alternating each 
other, followed by a new line.

The BODY can contain of multiline text.



## Known Limitations

- The OrbitDB that Superhighway84 uses is a public database, meaning everyone
  can alter its data. Since its using a standard _docstore_, PUT and DELETE
  events can alter existing data. This issue will be solved in the future by
  customizing the store to ignore these types of events.

- Superhighway84 is bound to the version of IPFS that Berty decides to support 
  for go-orbit-db. go-orbit-db updates, on the other hand, seem to introduce
  breaking changes from time to time, which are hard to debug as someone without
  in-depth knowledge nor documentation. Since Superhighway84 is pretty much a
  one-man-show it would be quite challenging to fork go-orbit-db in order to
  keep it up to date with IPFS and make its interface more stable. Unfortunately
  there doesn't seem to be an alternative to Berty's go-orbit-db as of right
  now, so Superhighway84 is basically stuck with it.
  If you happen to know your way around IPFS and maybe even go-orbit-db, and
  would like to support this project, please get in touch!

- If you have a newer IPFS version installed than the one used by
  Superhighway84, please make sure to **not upgrade** the IPFS_REPO that
  Superhighway84 is using. Otherwise you will get an error when starting
  Superhighway84 that will tell you that there is an IPFS repository mismatch:

  ```
  > panic: Your programs version (11) is lower than your repos (12).
  ```

  If this should be the case, please follow the instructions provided here:

  https://github.com/mrusme/superhighway84/issues/42#issuecomment-1100582472

- If you encounter the following issue your IPFS repo version might be older
  than what Superhighway84 is using:

  ```
  > panic: ipfs repo needs migration
  ```

  In this case you might want to follow the IPFS migration guide here:

  https://github.com/ipfs/fs-repo-migrations/blob/master/run.md

  Alternatively use the same IPFS version as used by Superhighway84 to
  initialize a dedicated Superhighway84 repository. Please refer to the
  INSTALLATION part for how to do so.



## Credits

- Superhighway84 name, code and graphics by [mrusme](https://github.com/mrusme)
- Logo backdrop by
  [Swift](https://twitter.com/Swift_1_2/status/1114865117533888512)


