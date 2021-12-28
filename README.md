Very crude Windows port. Confirmed working on Windows 10 Pro N 64 bit. Can read, reply and post new thread.

How to Install :

1. Install IPFS on Windows, add IPFS on PATH.
2. Create folder "b"on D: hard drive (or, setup your own preferred path on superhighway84.go line 42). This folder will store the database.
3. Customize your own username by modifying /tui/mainscreen.go on line 280 and line 318 (newArticle.From = "your username")
4. Build the binary (go build .)

Or, just run the precompiled binary. Make sure the IPFS is installed and D:/b folder is exist. 
