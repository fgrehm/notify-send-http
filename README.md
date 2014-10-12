# notify-send-http

Trigger `notify-send` across the network using HTTP, useful for triggering
notifications from local VMs / Containers into your own computer. It even supports
notification icons!

![demo](http://i.imgur.com/51hGcuW.png)

**_Tested on Ubuntu 14.04 only_**

## Why?

Because I do all of my dev work on virtualized environments and I use [guard](https://github.com/guard/guard/)
quite a lot to keep my Ruby tests running when files get changed. The problem is
that its builtin notification support will trigger a `notify-send` inside the virtual
environment instead of my machine.

With `notify-send-http` I can run an HTTP server on my machine and make use of a
custom `notify-send` executable on my virtual environments that has the same
interface as the original command and will send notifications to the HTTP server
so that I can see alerts poping up on my screen whenever the build fails.

## Server setup

First you'll need to download the HTTP server on the machine where notifications
will get displayed and drop it somewhere on your `PATH`.

For example, if `$HOME/bin` is on your `PATH`:

```sh
curl -L https://github.com/fgrehm/notify-send-http/releases/download/v0.1.0/server > $HOME/bin/notify-send-server
chmod +x $HOME/bin/notify-send-server
```

Then fire up the server with:

```sh
PORT=12345 notify-send-server
```

If you are on an Ubuntu machine, you can set it up to start automatically after
logging in to the machine:

```sh
cat <<-STR > ~/.config/autostart/notify-send-server.desktop
[Desktop Entry]
Type=Application
Exec=env PORT=12345 ${HOME}/bin/notify-send-server
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name[en_US]=NotifySend server
Name=NotifySend server
Comment[en_US]=
Comment=
STR
```

## Client setup

From another machine / VM / Linux Container you'll have to first download the
client somewhere on the `PATH`:

```sh
curl -L https://github.com/fgrehm/notify-send-http/releases/download/v0.1.0/client | sudo tee /usr/local/bin/notify-send &>/dev/null
sudo chmod +x /usr/local/bin/notify-send
```

Then point the CLI to the notification server:

```sh
SERVER_IP=$(ip route|awk '/default/ { print $3 }')
export NOTIFY_SEND_URL="http://${SERVER_IP}:12345"
```

And trigger notifications with:

```sh
notify-send "A summary" "Some message"
```

## Integration

### [Devstep](http://fgrehm.viewdocs.io/devstep/)

Download the client to a folder under `$HOME/devstep/bin`:

```sh
DEVSTEP_DIR="$HOME/devstep/bin"
mkdir -p $DEVSTEP_DIR
curl -sL https://github.com/fgrehm/notify-send-http/releases/download/v0.1.0/client > $DEVSTEP_DIR/notify-send
chmod +x $DEVSTEP_DIR/notify-send
```

Add to your `$HOME/devstep.yml`:

```yaml
volumes:
  # Share executable with container
  - '{{env "HOME"}}/devstep/bin/notify-send:/.devstep/bin/notify-send'
environment:
  # Set to a different IP if your docker0 bridge is set to something else
  NOTIFY_SEND_URL: "http://172.17.42.1:12345"
```

### [Docker](https://www.docker.com/)

```sh
# Grab Docker bridge IP
DOCKER_BRIDGE_IP=$(/sbin/ifconfig docker0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')

# Download binary
curl -sL https://github.com/fgrehm/notify-send-http/releases/download/v0.1.0/client > /tmp/notify-send && chmod +x /tmp/notify-send

# Start container ready to send notifications
docker run -ti --rm \
           -e NOTIFY_SEND_URL="http://${DOCKER_BRIDGE_IP}:12345" \
           -v /tmp/notify-send:/usr/bin/notify-send \
           -v `pwd`/success.png:/tmp/success.png \
           ubuntu:14.04 \
           /usr/bin/notify-send "Hello docker" "It Works!" -i /tmp/success.png
```

### [Vagrant](http://www.vagrantup.com/)

```ruby
# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
Vagrant.configure("2") do |config|
  config.vm.box = "<YOUR BOX>"
  config.vm.provision :shell,
                      path: 'https://github.com/fgrehm/notify-send-http/raw/master/vagrant-installer.sh',
                      args: ['12345']
end
```

## TODO

- [ ] Handle other `notify-send` parameters
- [ ] Implement support for other notifiers / platforms
