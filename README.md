# Glow - A CLI Tool to adapt Git-Flow

![glow logo](./assets/glow-logo.svg)

## Installation

**Mac OS**

```bash
brew tap meinto/glow https://github.com/meinto/glow
brew install meinto/glow/glow
```

**Linux**

```bash
GVERSION=$(curl https://api.github.com/repos/meinto/glow/releases/latest -s | jq .name -r | cut -c 2-)

# download i386 architecture
curl -Lo glow.tar.gz https://github.com/meinto/glow/releases/download/v${GVERSION}/glow_${GVERSION}_linux_i386.tar.gz
# or x86_64 architecture
curl -Lo glow.tar.gz https://github.com/meinto/glow/releases/download/v${GVERSION}/glow_${GVERSION}_linux_x86_64.tar.gz

tar -xvzf glow.tar.gz glow
mv glow /usr/local/bin/git-glow
rm glow.tar.gz
```

**manually**

Here you can find all [available Binaries](https://github.com/meinto/glow/releases). Download the binary and run the install command:

```bash
<name-of-binary> install
```

## Workflow

> Important!  
> Some commands need additional information like git author or Gitlab namespace and project name.  
> These informations can be stored in a config file or can be passed through flags.
> To configure the glow.json config file run the "init" command

![glow workflow](./assets/glow.jpg?raw=true)

### Feature Development

The following command will create a new feature branch in the following form: `features/dvader/death-star`. The name of the author (`dvader`) is grabbed from the config file.

```bash
# author grabbed from config
glow feature death-star
```

After you created the feature branch it is automatically checked out.  
When you finish your feature you can create a merge request in Gitlab:

```bash
# Gitlab information grabbed from config
glow close
```

### Create a release

I recommend to use [Semver](https://semver.org/) for versioning. The following command will create a release branch with the following format: `release/v1.2.3`.

```bash
glow release 1.2.3
```

### Publish a release

When you decide that the release is stable and you want to publish it, the following command will create a merge request on the `master` branch in Gitlab.

```bash
glow publish
```

### Close a release

After publishing the release, you have to merge all changes made on the release branch back into `develop`. The following command creates a merge request of the release branch into `develop`.

```bash
glow close
```

## Config

For some commands you must provide information like the url of your Gitlab instance or your Gitlab ci token. These informations can be put in a `glow.json` file. Glow will lookup this json in the directory where its executed.

You can create this json with the `init` command. The json will be automatically added to the `.gitignore`:

```bash
glow init
```

**List of all config params**

*glow.config.json*

```json
{
  "author": "dvadar",
  "gitProviderDomain": "https://gitlab.com",
  "gitProvider": "gitlab",
  "projectNamespace": "my-namespace",
  "projectName": "my-project",
  "gitPath": "/usr/local/bin/git",
  "mergeRequest.squashCommits": true,
  "versionFile": "VERSION",
  "versionFileType": "raw"
}
```

*glow.private.json*

```json
{
  "token": "abc"
}
```

## Git bindings

glow uses the native git installation per default. The default configured path to git is `/usr/local/bin/git`. You can change the path with the flag `--gitPath` or the property `gitPath` in the config file.
