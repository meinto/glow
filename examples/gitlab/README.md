# Example Gitlab ci pipeline

Here you can see how you could integrate `glow` in your ci/cd pipeline in Gitlab. The pipeline suggestion contains a bunch of manual triggers which create release branches and create merge requests for features, fixes, hotfixes and releases.

For the semantiv version number handling i recommend to use [git-semver](https://github.com/meinto/git-semver).