#!/bin/bash

# increment version
# commit the new versin changes
# tag the commit with the new version
# push all changes including tags
git semver version $1 -TP