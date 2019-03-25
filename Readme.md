## Git commit message formatter

Cli app to implement git message template  
Automate your commit message workflow with some rules

The name 'gitcomm' is short 'git commit message'.

Usage:

![demo](https://github.com/karantin2020/gitcomm/raw/master/docs/screen.gif)

```sh
Usage: gitcomm [-Avs]

Automate git commit messaging

Options:
  -V, --version   Show the version and exit
  -A, --addAll    Adds, modifies, and removes index entries to match the working tree. Evals `git add -A`
  -v, --verbose   Switch log output
  -s, --show      Show last commit or not. Evals `git show -s` in the end of execution
```

```sh
# type(<scope>): <Subject>

# <Body>

# * <Footer>
```

#### Type should be one of the following:

- feat (new feature)
- fix (bug fix)
- docs (changes to documentation)
- style (formatting, missing semi colons, etc; no code change)
- refactor (refactoring production code)
- test (adding missing tests, refactoring tests; no production code change)
- chore (updating grunt tasks etc; no production code change)
- version (description of version upgrade)

Scope is just the scope of the change. Something like (admin) or (teacher).  
Subject should use imperative tone and say what you did.  
The body should go into detail about changes made. Use the body to explain what and why vs. how  
The footer should contain any git issue references or actions.

#### Issue Processing

ISSUE_KEY #comment This is a comment  
ISSUE_KEY #resolved

Template is from this repo https://github.com/williampeterpaul/.git-commit-template

#### Follows this https://chris.beams.io/posts/git-commit/

#### The seven rules of a great Git commit message

> Keep in mind: This has all been said before.

1. Separate subject from body with a blank line
2. Limit the subject line to 72 characters
3. Capitalize the subject line
4. Do not end the subject line with a period
5. ~~Use the imperative mood in the subject line~~ // That rule is for user impl
6. Wrap the body at 320 characters
7. ~~Use the body to explain what and why vs. how~~ // That rule is for user impl

Example commit message:

    Redirect user to the requested page after login

    https://trello.com/path/to/relevant/card

    Users were being redirected to the home page after login, which is less
    useful than redirecting to the page they had originally requested before
    being redirected to the login form.

    * Store requested path in a session variable
    * Redirect to the stored location after successfully logging in the user

Source is https://robots.thoughtbot.com/5-useful-tips-for-a-better-commit-message
