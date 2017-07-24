# gorl

After losing interest in it for a while, the project which has eaten at me since
I was very young has returned to my attention again.

When I was 12 or 13 I had my first major programming project: I hacked together
a roguelike game in GameMaker. The code was a mess, the ideas were all over the
place, the enemy AI was dumb... *but it fucking worked*. Ever since that day,
I've been trying to recreate the magic. Embarrassingly, this is my fifth
attempt; after 2 in C, one in Perl, and one in C++, I think I've finally found
the language which I can use for this.

The details of the setting, theme, etc. will be left for later. Right now, I'm
focused on hacking together something worth playing.

## Code structure

The general idea is to separate the game logic from the graphics/input. gorl is
where all the common code goes, but it alone will not build an exectuable.
Different graphics backends will be implemented as separate executables in the
`cmd` directory. So to build the termbox version, you do:

``` shellsession
$ go get -u github.com/japanoise/gorl/cmd/gorl-termbox
```

The interfaces needed are in `interfaces.go`; sprites needed in `sprite.go`, and
tiles needed are in `tiles.go`

## Licensing

Code (.go files) Copyright © chameleon 2016-2017, licensed under the MIT license.
