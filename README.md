CloudFoundry Fast Push Plugin
==

`cf fast-push <app>` enables developers to push incremental updates of your application to CloudFoundry. Deploy in seconds rather than minutes today! This is possible by modifying your deployment configuration. For example your `manifest.yml` has to be changed to kick off the [cf-fastpush-controller](https://github.com/xiwenc/cf-fastpush-controller) instead of your application. Our [example python application](https://github.com/xiwenc/cf-python-test-app-fastpush) can be used to quickly test drive it.

During CloudFoundry Summit Berlin 2015 [Jouke](https://github.com/jtwaleson) and [Xiwen](https://github.com/xiwenc) had an idea how to speed up deployment speed of CF apps. The original idea was first brought up by Jouke because we were challenged to make application deployments faster at [Mendix](https://www.mendix.com). At the conference we were inspired by the community and also [Cloud Rocker](https://github.com/CloudCredo/cloudrocker). With `cf fast-push` we solve somewhat the same problem as `Cloud Rocker` namely: how can we develop applications at a faster pace. Cloud Rocker enables delopers to run their application locally using `Docker`. `fast-push` shortens deployment time to a real CloudFoundry cluster by efficiently synchronize changed files without restaging the application. This way we can still leverage all the goodies of CF like publicly accessible and the rich services ecosystem.

Current documentation is fairly low because this project is still very young. We invite others to contribute. The code quality is not great either and we plan to refactor it in the near future. Please keep in mind the current state is far from production ready but for its purpose it is good enough to be used ;)

Tutorial
===

First of all you need to install this plugin:
```bash
git clone https://github.com/xiwenc/cf-fastpush-plugin.git
cd cf-fastpush-plugin
go build
cf install-plugin cf-fastpush-plugin
```

Your application (tested with python so far but should work with any types that are supported by CloudFoundry) needs to be changed. As a requirement we need to include the executable of [cf-fastpush-controller](https://github.com/xiwenc/cf-fastpush-controller):
```bash
git clone https://github.com/xiwenc/cf-fastpush-controller.git
cd cf-fastpush-controller
go build
cp cf-fastpush-controller /path/to/my/app/cf-fastpush-controller
```

A `cf-fastpush-controller.yml` needs to be created with at least the command that needs to be executed to run your application. Example `cf-fastpush-controller`:
```yaml
backend_command: python hello.py
```
All the options that can be overridden can be found in [contants.go](https://github.com/xiwenc/cf-fastpush-controller/blob/master/lib/constants.go) and the defaults in [main.go](https://github.com/xiwenc/cf-fastpush-controller/blob/master/main.go).

Next the `command` field in your `manifest.yml` needs to be changed to start `cf-fastpush-controller` instead:
```yaml
command: ./cf-fastpush-controller
```

With these in place we can push the initial application with:
```bash
cf push my-cool-app
```

And incremental updates can be pushed:
```bash
cf fp my-cool-app
```

How it works
===

The fastpush mechanism uses the server-client model. The server is `cf-fastpush-controller` and the client is a cf cli plugin `cf-fastpush-plugin`.

- `cf-fastpush-controller`: A daemon that sits between your application and the gorouters. This service is always available and it responds to some specific paths under `/_fastpush/*`. Paths that are not known to the controller are reverse-proxied to the backend application which is your application. It keeps track of your remote files and accepts new and existing files. Depending on what files has changed it can trigger an automatic restart of the backend.
- `cf-fastpush-plugin`: A cf cli plugin that talks to the controller. It tracks your local files and synchronizes those that has been changed or added.

The actual code does more than what is documented here. So we suggest you to read the source if you are really interested in how it works and what else it can do.

Credits:
===
- Mendix: Great place to work
- Colleagues @ Mendix
