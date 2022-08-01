# Git Manager
[![Version](https://img.shields.io/badge/Version-0.2.0-blue)](https://gitlab.com/DeveloperC/git-manager/-/releases)
[![Pipeline Status](https://gitlab.com/DeveloperC/git-manager/badges/main/pipeline.svg)](https://gitlab.com/DeveloperC/git-manager/-/pipelines)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![License](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

// It sounds like you're describing what `git branch --set-upstream-to=$otherBranch` does?
// With `git rebase --onto <new_base> <previous_base>` you can achieve this, but it isn't as easy as it could be.

#### Flow
// Is a rebase in progress?
// If it is just warn and do not touch repo.

// Gather information.
		// Get origin/HEAD branch name.
		// git name-rev --name-only origin/HEAD

		// Get the branch currently checked out.
		// git branch --show-current

		// Get if their are uncommited changes.
		// git diff --name-only

// If uncommited changes on not mainline.
// git commit -m "temp"

// Update from remote.
		// Switch to mainline
		// git checkout $(git name-rev --name-only origin/HEAD)

		// Prune deletes merged branches etc.
		// git fetch --prune --prune-tags

		// Pull from remote. //TODO reset instead?
		// git", "pull", "--rebase", "--autostash

// For each local branch.
// git branch --format "%(refname:short) %(upstream:short)"
		// Rebase upon the updated head.
		// git checkout ${branch} && git rebase ...

		// Do we have unpushed branches?
		// git puch ${unpushed branch}

// Do we have branches with unpushed commits?
// ???
		// Will a simple push work?
		// git push ${branch}
