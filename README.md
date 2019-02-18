# Glow - Git-Flow for Gitlab

## Commands

### Create a feature branch

```bash
# this will create a feature branch like:
# features/dvader/death-star
glow feature death-star --author dvader

# shorter when you provide your author name in the glow.json file
glow feature death-star
```

### Create a release branch

```bash
# this will create a release branch like:
# release/v1.2.3
glow release 1.2.3

# When you want to do some actions after createing the release branch, for example to increase the version of your product, you can provide a post release script
glow release 1.2.3 --postRelease increaseVersion.sh
```

### Create a merge request

```bash
# I recommend to provide the properties
# - gitlabEndpoint
# - projectNamespace
# - projectName
# - gitlabCIToken
# in the glow.json.
glow mergeRequest source/branch target/branch

# If you don't want to use the config file you can do ist all on the command line:
glow mergeRequest source/branch target/branch \
  -e https://gitlab.com \ # gitlabEndpoint
  -n my-namespace \       # projectNamespace
  -p my-project \         # projectName
  -t abc                  # gitlabCIToken
```

## Config

For some commands you must provide information like the url of your gitlab instance or your gitlab ci token. These informations can be put in a `glow.json` file. Glow will lookup this json in the directory where its executed.

**List of all config params**

```json
{
  "author": "dvadar",
  "gitlabEndpoint": "https://gitlab.com",
  "projectNamespace": "my-namespace",
  "projectName": "my-project",
  "gitlabCIToken": "abc",
}
```