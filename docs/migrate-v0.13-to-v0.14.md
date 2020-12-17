# Migration Guide: v0.13 to v0.14

v0.14.0 introduces a number of breaking changes, throu it should not be hard to migrate.  
Just follow this guid and if issues still ocure ask for help on discord or  
feel free to create an issue.

<!-- toc -->

-   [Removed Functions (#467)](#Removed-Functions)
-   [Renamed Functions (#467)](#Renamed-Functions)

<!-- tocstop -->

## Removed Functions

 - for **GetUserTrackedTimes** and **GetRepoTrackedTimes** use **ListRepoTrackedTimes** with specific options set

Pulls:
-   [#467 Remove GetUserTrackedTimes](https://gitea.com/gitea/go-sdk/pulls/467)


## Renamed Functions

- **ListTrackedTimes** is now **ListIssueTrackedTimes**

Pulls:
-   [#467 Remove & Rename TrackedTimes list functions](https://gitea.com/gitea/go-sdk/pulls/467)
