# Migration Guide: v0.11 to v0.12

v0.12.0 introduces a number of breaking changes, throu it should not be hard to
migrate.
Just follow this guid and if issues still ocure ask for help on discord or
feel free to create an issue.

<!-- toc -->

-   [List Functions now always need an ListOption as argument (#205) (#243) (244)](#List-Functions-now-always-need-an-ListOption-as-argument)
-   [Authentification was removed from all Functions (#241) (#335)](Authentification-was-removed-from-all-Functions)
-   [Some Functions where deleted (#247) (#262)](Some-Functions-where-deleted)
-   [SearchUsers arguments are move to an Option struct (#248)](SearchUsers-arguments-are-move-to-an-Option-struct)
-   [RepoTopics functions now expect and return string slice directly (#276)](ListRepoTopics-return-now-string-slice-directly)
-   [MergePullRequestOption field names changed and Enum is now used (#328)](MergePullRequestOption-field-names-changed-and-Enum-is-now-used)

<!-- tocstop -->

## List Functions now always need an ListOption as argument

since paggination is introduced in gitea v1.12.0 for all list endpoints,
all List Functions acept at least **Page** and **PageSize**.

If the function had had already an Option struct as argument this one is now extendet,
if not a new Options type was created.

-   migrate old paggination arguments to the new One.
-   add a empty Option struct if a new one was created.

Pulls:

-   [#205 Add Pagination Options for List Requests](https://gitea.com/gitea/go-sdk/pulls/205)
-   [#243 Add ListIssueCommentOptions for optional param](https://gitea.com/gitea/go-sdk/pulls/243)
-   [#244 Add ListMilestoneOption to ListRepoMilestones](https://gitea.com/gitea/go-sdk/pulls/244)

## Authentification was removed from all Functions

for Authentification the default credentials/token is used,
witch was set on Client initialisation.

for RepoWatch functions remove arguments:

-   GetWatchedRepos: password (second)
-   WatchRepo: username (first), password (second)
-   UnWatchRepo: username (first), password (second)

for Token functions remove:

-   the first two argument (user & password),
    these functions still relay on BasicAuth so if not done,
    just set username, password and optional otp before executing them.

```go
client.SetBasicAuth(username, password)
client.SetOTP(otp)
```

Pulls:

-   [#241 Refactor RepoWatch](https://gitea.com/gitea/go-sdk/pulls/241)
-   [#335 Support 2FA for basic auth & Refactor Token functions](https://gitea.com/gitea/go-sdk/pulls/335)

## Some Functions where deleted

Functions where deleted because they where only workarounds
or are helper functions witch could be replaced easely.

-   BasicAuthEncode
    if you realy need this just copy the function into your project:
    ```go
    func BasicAuthEncode(user, pass string) string {
    		  return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
    }
    ```
-   ListUserIssues
    If you realy need this just use the Workaround witch was removed with #262
    and If you have time a pull upstream to gitea for a real API is always wellcome

Pulls:

-   [#247 Fix & Refactor UserApp Functions](https://gitea.com/gitea/go-sdk/pulls/247)
-   [#262 Remove ListUserIssues](https://gitea.com/gitea/go-sdk/pulls/262)

## SearchUsers arguments are move to an Option struct

Old:
 `client.SearchUsers(query, limit)`
New:
 `client.SearchUsers(SearchUsersOption{KeyWord: "query", Page: 1, PageSize: limit})`

Pull: [#248 extend SearchUsers](https://gitea.com/gitea/go-sdk/pulls/248)

## ListRepoTopics return now string slice directly

ListRepoTopics returned a struct with Topics string slice.
Now it return the falue of this string slice directly

Old:

```go
client.SetRepoTopics(user, repo, TopicsList{topic_slice})
```

New:

```go
client.SetRepoTopics(user, repo, topic_slice)
```

Pull: [#276 Refactor List/SetRepoTopics](https://gitea.com/gitea/go-sdk/pulls/276)

## MergePullRequestOption field names changed and Enum is now used

Rename **MergeTitleField** to **Title**
Rename **MergeMessageField** to **Message**

Do is now called Style and expect predefined falues:
MergeStyleMerge, MergeStyleRebase, MergeStyleRebaseMerge & MergeStyleSquash

Pull: [#328 PullMerge: use enum for MergeStyle](https://gitea.com/gitea/go-sdk/pulls/328)
