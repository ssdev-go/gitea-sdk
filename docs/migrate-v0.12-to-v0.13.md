# Migration Guide: v0.12 to v0.13

v0.13.0 introduces a number of breaking changes, throu it should not be hard to migrate.  
Just follow this guid and if issues still ocure ask for help on discord or  
feel free to create an issue.

<!-- toc -->

-   [EditMilestoneOption use StateType (#350)](#EditMilestoneOption-use-StateType)
-   [RepoSearch Options Struct was rewritten (#346)](#RepoSearch-Options-Struct-was-rewritten)
-   [Variable Renames (#386)](#Variable-Renames)
-   [Change Type of Permission Field (#408)](#Change-Type-of-Permission-Field)
-   [All Function return http responce (#416)](#All-Function-return-http-responce)

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

## Change Type of Permission Field

The following functions are affected:  ListOrgTeams, ListMyTeams, GetTeam, CreateTeam, EditTeam and AddCollaborator

The `Permission` field has changed type from `string` to `AccessMode`,  
which represent the raw strings you must use before.  
Just replace the string with the AccessMode equivalent.

Pulls:
-   [#408 Use enum AccessMode for OrgTeam and Collaborator functions](https://gitea.com/gitea/go-sdk/pulls/408)


## All Function return http responce

All functions got one new return (`Responce`)!  
If you just like to migrate, add `_,` before the error return.

example:
```diff
- user, err := c.GetMyUserInfo()
+ user, _, err := c.GetMyUserInfo()
```

If you like to check responce if an error ocure, make sure responce is not nil!  
If an error ocure before an http request (e.g. gitea is to old), it will be nil.

Pulls:
-   [#416 All Function return http responce](https://gitea.com/gitea/go-sdk/pulls/416)
