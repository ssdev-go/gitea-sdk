# Migration Guide: v0.13 to v0.14

v0.14.0 introduces a number of breaking changes, throu it should not be hard to migrate.  
Just follow this guid and if issues still ocure ask for help on discord or  
feel free to create an issue.

<!-- toc -->

-   [Removed Functions (#467)](#removed-functions)
-   [Renamed Functions (#467)](#renamed-functions)
-   [New Optional Fields (#486)](#new-optional-fields)

<!-- tocstop -->

## Removed Functions

 - for **GetUserTrackedTimes** and **GetRepoTrackedTimes** use **ListRepoTrackedTimes** with specific options set

Pulls:
-   [#467 Remove GetUserTrackedTimes](https://gitea.com/gitea/go-sdk/pulls/467)


## Renamed Functions

- **ListTrackedTimes** is now **ListIssueTrackedTimes**

Pulls:
-   [#467 Remove & Rename TrackedTimes list functions](https://gitea.com/gitea/go-sdk/pulls/467)


## New Optional Fields

The `EditUserOption` struct has gained several new Optional fields.
For example Email type changed from `string` to `*string`.

The easiest migration path is, to wrap your options with:
**OptionalString()**, **OptionalBool()** and **OptionalInt64()**

Pulls:
-   [#486 Update Structs](https://gitea.com/gitea/go-sdk/pulls/486)
