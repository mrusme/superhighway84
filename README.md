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

```

INSTALLATION
------------

Clone this repository and run:

$ go build .

The binary will be available at ./superhighway84 and can be moved wherever you
please.

If you don't have IPFS installed already, make sure to do so in order to be able
to initialize your IPFS repository:

https://docs.ipfs.io/install/command-line/

The IPFS repository can be initialized using the following command:

$ ipfs init



RUNNING
-------

First, check ulimit -n and verify that it's at a reasonable amount. IPFS
requires it to be large enough (>= 2048) in order to work properly over time.

Second, if your hardware shouldn't be a beefy computer but instead one of
those flimsy MacBooks, older hardware, a Raspberry or a low-memory VPS it is
advisable to set the previously created IPFS repository to the `lowpower`
profile.

$ ipfs config profile apply lowpower

This should help with CPU usage, file descriptors and the amount of network
connections. While during the startup period you might still see peers peaking
between 1k and 3k, connections should ultimately settle somewhere between 100
and 300 peers.

Afterwards you can simply launch the binary:

$ ./superhighway84

A setup wizard will help you with initial configuration. Please make sure to
have at least HOME and EDITOR exported in your environment.

In case you would like to use a dedicated ipfs repository for Superhighway84,
you will have to export a different IPFS_PATH and make sure it was initialized
beforehand:

$ export IPFS_PATH=~/.ipfs-sh84
$ ipfs init
$ superhighway84

In case you're intending to run the official IPFS daemon and Superhighway84 in
parallel, be sure to adjust the ports in their respective IPFS repos (e.g.
~/.ipfs and ~/.ipfs-sh84) so that they won't utilize the same port numbers.
The ports 4001, 5001 and 8080 are relevant and should be adjusted to something
other for every new repo/IPFS node that will run in parallel, e.g.:

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

NOTE: When running Superhighway84 for the first time it might seem like it's
"hanging" at the command prompt. Usually it isn't hanging but rather searching
for peer it can connect to in order to synchronize the database. Depending on
how many people are online, this process might take _some time_, please be
patient.



CONFIGURATION
-------------

Superhighway84 will guide you through the basic configuration on its first run.
The configuration is stored at the path that you specified in the setup wizard.
After it was successfully created, it can be adjusted manually and will take
effect on the next launch of Superhighway84.

Configuration options that might be of interest:

- superhighway84.toml - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

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

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -



USAGE
-----

The default keyboard shortcuts are:

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

However, you are free to customize these within your configuration file, under
the section `Shortcuts`. 


SUBMIT ARTICLE

When submitting a new article or a reply to an article, the $EDITOR is launched
in which a document with a specific structure will be visible. This structure
consists of the HEADER, a SEPARATOR and the BODY and looks like this:

- $EDITOR - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

Subject: This is the subject of the article
Newsgroup: test.sandbox
= = = = = =
This is the multiline
body of the article

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

The HEADER contains all headers that are required for an article to be
submitted. These are:

- Subject: 
  The subject of the article that will be shown in the articles list. The
  subject must only contain of printable ASCII characters.

- Newsgroup: 
  The newsgroup under which the article will be submitted, this can
  either be an existing group or a new group. Please try to follow
  the convention when creating new groups.
  The newsgroup must only contain of printable ASCII characters.

The SEPARATOR contains of 6 equal signs and 5 spaces, alternating each 
other, followed by a new line.

The BODY can contain of multiline text.



KNOWN LIMITATIONS
-----------------

- The OrbitDB that Superhighway84 uses is a public database, meaning everyone
  can alter its data. Since its using a standard _docstore_, PUT and DELETE
  events can alter existing data. This issue will be solved in the future by
  customizing the store to ignore these types of events.
- Superhighway84 is always behind recent IPFS and also OrbitDB versions, mainly
  because Berty, the go-orbit-db maintainers, aren't exactly super helpful 
  and welcoming in regard of the usage of their library. Not only do they
  simply not document a thing or take interest in answering issue reports on
  GitHub, they also don't seem to care about supporting recent IPFS versions 
  either. 
  Superhighway84 is bound to the version of IPFS that Berty decides to support 
  for go-orbit-db. go-orbit-db updates, on the other hand, seem to introduce
  breaking changes from time to time, which are hard to debug as someone without
  in-depth knowledge nor documentation, and get basically no support from the
  Berty developers whatsoever. Since Superhighway84 is pretty much a
  one-man-show it would be quite challenging to fork go-orbit-db in order to
  keep it up to date with IPFS and make its interface more stable. Unfortunately
  there doesn't seem to be an alternative to Berty's go-orbit-db as of right
  now, so Superhighway84 is basically stuck with it.
  If you happen to know your way around IPFS and maybe even go-orbit-db, and
  would like to support this project, please get in touch!



CREDITS
-------

- Superhighway84 name, code and graphics by mrusme
  https://github.com/mrusme

- Logo backdrop by Swift
  https://twitter.com/Swift_1_2/status/1114865117533888512

```

