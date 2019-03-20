# Example Gitlab ci pipeline

Here you can see how you could integrate `glow` in your ci/cd pipeline in Gitlab. The pipeline suggestion contains a bunch of manual triggers which create release branches and create merge requests for features, fixes, hotfixes and releases.

For the semantiv version number handling [git-semver](https://github.com/meinto/git-semver) is required by this library.

## The glow stage

All actions which are represented by manual triggers in the Gitlab pipeline take place at the `glow` stage.

```yml
stage:
  - test
  - build
  - glow    # <-- this one
```

Depending on the conditions defined with the `only` property of a Gitlab ci job definition, the corresponding job is only present on specific branches for example.

## Creating a release

![create release](./create-release.png)

The following example shows how to trigger a release branch creation of a `patch` release. You could copy this section and addapt it for `minor` and `major` releases of course.

```yml
create:patch:release:
  stage: glow
  script: 
    - git glow cleanup branches --untracked --gone
    - git glow cleanup tags --untracked
    - git glow release patch
        --push
        --cicdOrigin https://$CI_GIT_USER:$CI_GIT_TOKEN@gitlab.com/$NAMESPACE/$PROJECT/
  only:
    - develop
  when: manual
```

First of all we cleanup the repository managed by the Gitlab runner with the first two cleanup scripts `git glow cleanup ...`.

After that the `release` command would create a patch release branch, and push it to the bare repository. 

To make this possible, you have to make shure, that your Gitlab runner have access rights to push to your repository. One way to do so is to set a new origin with an access token. In the example above you can see how this would look like.

> Of course this is not the most beautiful method. A better way would be to setup a seperate Gitlab ci user and allow this user accessing the repository. On your server you could configure your Gitlab user to access gitlab via ssh key.

## Publish a release

![publish release](./publish-release.png)

When you are shure that your release on your release branch is stable, trigger a publish merge request in Gitlab:

```yml
trigger:prod:release:
  stage: glow
  script:
    - git glow publish --ci
  only:
    - /^release\/v.*$/
    - /^hotfix\/.*$/
  when: manual
```

This manual trigger would create a merge-request in Gitlab to merge the feature with the `master` branch.

This trigger also applies to hotfix branches.

## Close a release

After creating a merge request on `master`, you should define another pipeline which creates another merge request of the release branch back into your develop branch, and tags your published release in git.

```yml 
trigger:close:release:
  stage: glow
  script:
    - git glow close release current
    - git glow tag
        --cicdOrigin https://$CI_GIT_USER:$CI_GIT_TOKEN@git.intern.hornbach.de/apps/emm-intents.git/
  only:
    - master
```

This pipeline job runs when the release branch was successfully merged into the master branch, by accepting the merge-request created in the step before. The first command `git glow close release current` closes the release branch with the current version number of the master branch. The second command creates a git tag and pushes it up to the bare repository.

