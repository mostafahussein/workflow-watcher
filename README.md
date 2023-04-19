# Workflow Watcher

[![ci](https://github.com/mostafahussein/workflow-watcher/actions/workflows/build.yaml/badge.svg)](https://github.com/mostafahussein/workflow-watcher/actions/workflows/build.yaml)

Pause a GitHub Actions workflow and wait for another workflow to complete before continuing.

Sometimes, a commit can result in cache invalidation, such as updating application dependencies, and you want to apply this commit to multiple branches. To maintain the ***Build Once, Deploy Anywhere*** principle in such cases, you can either wait until a specific branch is built before promoting your artifact to the next environment, or configure your workflow to check if there is an existing workflow running for the same commit in case you reset the other branches to a specific branch that contains the desired commits.


The way this action works is the following:

1. Workflow comes to the `workflow-watcher` action.
2. `workflow-watcher` will check if there is a workflow already running for the specified commit.
3. If and once the previously detected workflow is completed successfully, the workflow will continue.
4. If the previously detected workflow failed for any reason, then the workflow will exit with a failed status.


## Usage

```yaml
steps:
  - uses: mostafahussein/workflow-watcher@v1.0.0
    if: ${{ github.ref != 'refs/heads/develop' }}
    with:
      secret: ${{ secrets.GH_TOKEN }}
      repository-name: ${{ github.repository }}
      repository-owner: ${{ github.repository_owner }}
      head-sha: ${{ github.sha }}
      base-branch: "develop"
      polling-interval: 60

```

- `head-sha` is the commit SHA that triggered the workflow. The value of this commit SHA depends on the event that triggered the workflow. For more information, see "[Events that trigger workflows.](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows)" For example, `ffac537e6cbbf934b08745a378932722df287a53`.
- `base-branch` is the branch that will be used as a source for your final artifact. For example, the testing branch will used the same artifact from the develop branch once the build is done, in this case the `base-branch` value should be `develop`
- `polling-interval` determines how often a poll occurs to check for a updates from Github API, by default it will be **30 seconds**, in case you need more time or having issues with Github rate limiting, you can set your own polling interval.

## Timeout

If you'd like to force a timeout of your workflow pause, you can specify `timeout-minutes` at either the [step](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idstepstimeout-minutes) level or the [job](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes) level.

For instance, if you want your workflow watcher step to timeout after an hour you could do the following:

```yaml
steps:
  - uses: mostafahussein/workflow-watcher@v1.0.0
    timeout-minutes: 60
    ...
```

## Limitations

* While the workflow is paused, it will still continue to consume a concurrent job allocation out of the [max concurrent jobs](https://docs.github.com/en/actions/learn-github-actions/usage-limits-billing-and-administration#usage-limits).
* A job (including a paused job) will be failed [after 6 hours](https://docs.github.com/en/actions/learn-github-actions/usage-limits-billing-and-administration#usage-limits).
* A paused job is still running compute/instance/virtual machine and will continue to incur costs.

## Development

### Running test code

To test out your code in an action, you need to build the image and push it to a different container registry repository. For instance, if I want to test some code I won't build the image with the main image repository. Prior to this, comment out the label binding the image to a repo:

```dockerfile
# LABEL org.opencontainers.image.source https://github.com/mostafahussein/workflow-watcher
```

Build the image:

```
$ VERSION=1.1.1-rc.1 make IMAGE_REPO=ghcr.io/mostafahussein/workflow-watcher-test build
```

*Note: The image version can be whatever you want, as this image wouldn't be pushed to production. It is only for testing.*

Push the image to your container registry:

```
$ VERSION=1.1.1-rc.1 make IMAGE_REPO=ghcr.io/mostafahussein/workflow-watcher-test push
```

To test out the image you will need to modify `action.yaml` so that it points to your new image that you're testing:

```yaml
  image: docker://ghcr.io/mostafahussein/workflow-watcher-test:1.1.0-rc.1
```

Then to test out the image, run a workflow specifying your dev branch:

```yaml
- name: Watch Workflow on Develop branch
  uses: your-github-user/workflow-watcher@your-dev-branch
  if: ${{ github.ref != 'refs/heads/develop' }}
  with:
    secret: ${{ secrets.GH_TOKEN }}
    repository-name: ${{ github.repository }}
    repository-owner: ${{ github.repository_owner }}
    head-sha: ${{ github.sha }}
    base-branch: "develop"

```

For `uses`, this should point to your repo and dev branch.

## Credits

- Author: [Mostafa Hussein](https://github.com/mostafahussein)
- Inspired by: [Manual Workflow Approval](https://github.com/trstringer/manual-approval)