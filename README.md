# screech
A music player in the terminal designed for termux

DISCLAIMER: Screech is still in its alpha build and is still unstable, so please do well to report any issues you may encounter.


## Installation 

Screech is designed to be extremely easy to setup, only takes a few steps 

*Note*: Screech requires Termux:Api in order to work because it makes use of a number of apis.

**1. Clone the repo**
```bash 
git clone https://github.com/lordryns/screech.git
```

**2. Enter the location**
```bash 
cd screech
```

**3. Run the setup script**
```bash 
bash setup.sh
```

This should try to add screech to path and setup Termux:Api (if it isn't already configured)

**4. Restart Termux and start screech**
```bash 
screech
```

That's it! 


## Note 

On your first launch, Screech will try to scan your device for music, this is a one time thing and will not happen again unless you manually try to rescan by doing `Ctrl+f`.


## Other stuff 
This is just the initial build and so it might be unstable and is lacking features. Anyway, to contact me outside of Github, try @lordryns on [X](https://x.com/lordryns) and [Discord](https://discord.com/users/1015382973052358657)
