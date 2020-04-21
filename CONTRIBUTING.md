# Contributing to Clamshell

This guide will walk you through contributing to Clamshell, and assumes you
have some experience with the game go, but might be relatively new to
programming or to the tools/technologies.

## Bugs, Features, Questions, Comments

If you have any bugs, feature requests, questions, or comments, please file
them on our Issue on our [github tracker](https://github.com/otrego/clamshell/issues)

## Getting Started on Go

First, if you're new to Go, the language, check out:

1.  [Installing Go](https://golang.org/doc/install)
2.  [The Go Tour](https://tour.golang.org/welcome/1)
3.  [How to Write Go Code](https://golang.org/doc/code.html)

If you're very new to Go, you might want to check out the
[otrego/experimental](https://github.com/otrego/experimental)
repository.

## Getting Started with Git

If you're new to Git / Github, check out:

1.  [The Git Book](https://git-scm.com/book/en/v2). In particular, I find these
    sections very helpful:
    1.  [Basic Branching & Merging](https://git-scm.com/book/en/v2/Git-Branching-Basic-Branching-and-Merging)
    2.  [Distributed Workflows](https://git-scm.com/book/en/v2/Distributed-Git-Distributed-Workflows#ch05-distributed-git). 
        Specifically, we use an [Integration-Manager workflow (Github)](https://git-scm.com/book/en/v2/Distributed-Git-Distributed-Workflows#wfdiag_b).
2.  [Getting Started with Github](https://help.github.com/en/github/getting-started-with-github)

## Setting Up Your Dev Environment

First, make sure you have
[two-factor auth](https://help.github.com/en/github/authenticating-to-github/securing-your-account-with-two-factor-authentication-2fa)
setup in Github. This is a requirement for working with Otrego.

We expect all developers, both external contributors and maintainers to
interact with Clamshell via a
[fork-and-pull-request model](https://help.github.com/en/github/getting-started-with-github/fork-a-repo).
So generally, developers will fork the Clamshell repository and then submit
pull requests to the primary `otrego/clamshell` repository.

In general, once you install the Go programming language, you should set
`GOPATH` to something that suits your tastes. I (kashomon) set `GOPATH` to be
`$HOME/inprogress/go`, but feel free to set it as you wish. For more about
setting GOPATH, check out the 
[golang wiki](https://github.com/golang/go/wiki/SettingGOPATH).

Assuming you have created a fork of Clamshell (see Development Model), then
create the relevant directories & clone the repo:

```shell
mkdir -p $GOPATH/src/github.com/otrego
git clone github.com/<USERNAME>/clamshell $GOPATH/src/github.com/otrego/clamshell
```

Set the upsteams appropriately:

```shell
cd $GOPATH/src/github.com/otrego/clamshell
git remote add upstream git@github.com:otrego/clamshell.git

# Check everything's set up correctly:
git remote -v
```

Then, make sure it builds!

```shell
cd $GOPATH/src/github.com/otrego/clamshell
go test ./...
```

## Change Workflow

Now that you've set up your dev environment, let's talk about git development
workflow. The general flow should be:

1.  Do development on branches in your fork.
2.  Make pull requests to the main repo `otrego/clamshell`
3.  Get your code reviewed by a team member.
4.  Once approved, submit the changes.

Despite it's perils, I've found that rebase is a little easier to work with and
understand than merge when doing development on the local copy of your fork.

Word of caution: **Avoid using merge and rebase together**. That way lies peril.

First, here's what my workflow looks like (which is quite similar to this
[Atlassian guide](https://www.atlassian.com/git/tutorials/git-forks-and-upstreams).

### Example Workflow (kashomon)

1.  Rebase on any upstream changes into master via merge.

    ```shell
    git checkout master
    git fetch upstream
    git rebase -i upstream master

    # update my fork's master branch.
    git push
    ```

2.  Do feature development on a branch

    ```shell
    git checkout -b somefeature
    # ... do some work
    git add -A .
    git commit -a
    git push origin somefeature
    ```

3.  (Optional) If much time has passed, make sure your local master & feature
    branches are updated, running through 1. and then rebasing on top of that.

    ```shell
    git checkout master
    git fetch upstream
    git rebase -i upstream master
    git push

    # Update feature branch
    git checkout somefeature
    git rebase -i master

    git push

    # If necessary, you might need to force-push the changes, depending on the
    # nature of the rebase changes.
    git push --force
    ```

4.  When your change is ready, use the Github UI to get code reviewed & merge
    the changes into the repository:

    1.  Perform pull request via Github UI from your repsitory + feature branch
        to Github.

    2.  Get code reviewed by team member in Github UI. If you are reviewing
        code, make sure to 'Approve' the change.

    3.  Sqaush into single commit via Github UI and merge into the repository
        (this should happen by default.).
