# Shin: Shell input method

Shin lets you run shell commands anywhere you can type text. Watch this:

![Screencast](https://user-images.githubusercontent.com/2702526/197121193-1ec1956b-06d4-45bd-b145-b4c4aa82cf67.gif)

**This works in almost every Linux application,** including Firefox, Chrome,
LibreOffice, all GTK and Qt applications, and many more. The possibilities
are endless as static editing and dynamic evaluation become one.

Shin stays completely out of your way until you press the magic key that
brings it to life. When not in use, Shin consumes *zero* CPU time and memory.
I despise software that wastes computing resources, and Shin is designed
accordingly.

Ready to give it a try?


## Installation

Installing Shin is a painless procedure on most Linux distributions.

### Step 1: Set up IBus

Shin is implemented as an [IBus](https://github.com/ibus/ibus) input engine.
Before Shin can be installed, IBus must be installed and configured.

**If you use GNOME, you can skip this step,** because the GNOME desktop comes
with a fully configured IBus installation out of the box.

On KDE, it usually suffices to install any IBus input method (e.g. Typing Booster)
from the KDE Software Center, after which you will be prompted to switch to
IBus as the system input method. After logging out and back in, you should
be ready to go.

For other desktop environments and window managers, the setup may involve
manually installing the IBus daemon, and setting environment variables.
See [this ArchWiki article](https://wiki.archlinux.org/title/IBus) for more
information.

**Before filing an issue of the type "Shin is not working!", please verify
that you have a functioning IBus setup.** You can do this by trying out any
other IBus input method, such as Typing Booster or Pinyin.

### Step 2: Install Shin

Make sure you have [Go](https://go.dev/), a C compiler, Make, and Git.
Then run:

```
git clone https://github.com/p-e-w/shin.git
cd shin
make
sudo make install
ibus restart
```

*__Note 1:__ Occasionally, `ibus restart` fails to actually restart IBus.
In that case, you have to restart the IBus daemon manually, or simply
restart your system.*

*__Note 2:__ The above commands assume that your IBus installation is located
in `/usr/share`. In the vast majority of cases, this is correct. If you have
installed IBus elsewhere, you must run both `make` and `make install` with
the `IBUS_INSTALL_DIR` variable set to the actual location, e.g.
`make IBUS_INSTALL_DIR=/usr/share`.*

### Step 3: Make Shin activatable

The Shin input engine is intended to be activated using a hotkey, and
automatically deactivates itself after a shell command has been entered.
This allows for normal text editing to continue without having to manually
switch back to the default input method.

For this reason, **it is not recommended to add Shin to the desktop's
input method switcher.**

Instead, simply bind the command `ibus engine shin` to a global keyboard
shortcut such as <kbd>Alt</kbd>+<kbd>Space</kbd>. On GNOME, this can be done
using *Settings > Keyboard > View and Customize Shortcuts > Custom Shortcuts*.
On KDE, *System Settings > Shortcuts > Custom Shortcuts* does the same thing.
For other desktop environments and window managers, see the appropriate
documentation.

That's it! Shin is now just a keypress away whenever you need it.


## The `shin/bin` directory

If you use Shin frequently, it's quite natural to want to define custom commands
for inserting commonly needed text. For example, you might like to have a
`sig` command that inserts a signature with your name and address for use in
emails and online discussions.

Of course, you could just create a shell script which prints that text, and place
it in any of the directories in the shell's search path. But then `sig` would
also be available in regular interactive shells, and you probably don't want to
pollute your global command namespace with Shin-specific commands. You may also
wish to override some of the standard commands in Shin, *but only* in Shin.

To solve these problems, Shin prepends the directory `$XDG_CONFIG_HOME/shin/bin`
(which usually expands to `~/.config/shin/bin`) to the shell's search path when
running commands. To define Shin-specific commands, simply create that directory,
and drop executable scripts with the desired names there. They will be available
in Shin, without affecting the behavior of the shell anywhere else.


## Security considerations

By design, Shin turns every text input on your system into a basic terminal
emulator with full shell access. This shouldn't be a problem under normal
circumstances, since anyone with unsupervised access to your computer can just
launch a terminal anyway, but in some situations, such as a system running in
"kiosk mode", extra caution may be warranted.

Most screen locking applications explicitly disable input methods, but if you
use a non-standard screen locker, you should verify that Shin cannot be
accessed from the lock screen inputs, because that would create a trivial way
to bypass the lock by running a shell command that kills the lock screen process.

To the best of my knowledge, web browsers do not allow JavaScript code to
synthesize the low-level input events needed to control an IME. I therefore
believe that Shin is safe to use in browser inputs, even on untrusted websites.
A carefully designed website might however use concealed inputs and fake
input overlays to trick you into thinking that you have typed something
different from what you actually did. This represents a rather low security
risk though, since the site still cannot *control* which commands are entered
and executed.


## Acknowledgments

Shin depends on the [xdg](https://github.com/adrg/xdg),
[go-sqlite3](https://github.com/mattn/go-sqlite3), and
[dbus](https://github.com/godbus/dbus) Go packages,
as well as on
[BambooEngine's fork](https://github.com/BambooEngine/goibus)
of the goibus package (the only fork with proper Wayland support).

This project was my first encounter with IBus, a system that in practice
is used almost exclusively for typing East Asian scripts. I was pleasantly
surprised to find that IBus is flexible enough to support a use case like
Shin, which it was obviously not intended for and yet able to accommodate.
My respect goes to the creators of IBus for their clean and versatile design.


## License

Copyright &copy; 2022  Philipp Emanuel Weidmann (<pew@worldwidemann.com>)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

**By contributing to this project, you agree to release your
contributions under the same license.**
