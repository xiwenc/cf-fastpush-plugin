CloudFoundry Fast Push Plugin
==

`cf-fastpush-plugin` is a Cloud Foundry CLI plugin that speaks _fastpush_ REST API.

Usage
===

`cf-fastpush-plugin` on its own is not very useful. A fastpush controller like [cf-fastpush-controller](https://github.com/xiwenc/cf-fastpush-controller) is needed. For usage and general documentation on `fastpush` please refer to the documentation at [fastpush](https://github.com/xiwenc/fastpush).

Requirements
===

- Cloudfoundry CLI 6.x or later

Installation
===

```bash
git clone https://github.com/xiwenc/cf-fastpush-plugin.git
cd cf-fastpush-plugin
go build
cf install-plugin cf-fastpush-plugin
```

Commands
===

| Command | Short-cut | Description |
| --- | --- | --- |
| `cf fast-push <app name>` | `cf fp <app name>` | Update application files and restart app if needed. |
| `cf fast-push-status <app name>` | `cf fps <app name>` | Get status of the app. |
