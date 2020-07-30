# Migration Guide: v0.12 to v0.13

v0.13.0 introduces a number of breaking changes, throu it should not be hard to
migrate.
Just follow this guid and if issues still ocure ask for help on discord or
feel free to create an issue.

<!-- toc -->

-   [EditMilestoneOption use StateType (#350)](#EditMilestoneOption-use-StateType)
-   [RepoSearch Options Struct was rewritten (#346)](#RepoSearch-Options-Struct-was-rewritten)
-   [Variable Renames (#386)](#Variable-Renames)

<!-- tocstop -->

## EditMilestoneOption use StateType

Instead of a raw string StateType is now used for State too.
just replace old strings with new enum.


Pulls:

-   [#350 EditMilestoneOption also use StateType](https://gitea.com/gitea/go-sdk/pulls/350)


## RepoSearch Options Struct was rewritten

Since the API itself is ugly and there was no nameconvention whats o ever.
You easely can pass the wrong options and dont get the result you want.

Now it is rewritten and translated for the API.
The easyest way to migrate is to look at who this function is used and rewritten that code block.

If there is a special edgecase you have you can pass a `RawQuery` to the API endpoint.

Pulls:

-   [#346 Refactor RepoSearch to be easy usable](https://gitea.com/gitea/go-sdk/pulls/346)


## Variable Renames

Some names of strcut options have been renamed to describe there function/usecase more precisely.
if you use `CreateOrgOption` somewhere just rename `UserName` to `Name`.

Pulls:

-   [#386 CreateOrgOption rename UserName to Name](https://gitea.com/gitea/go-sdk/pulls/386)
