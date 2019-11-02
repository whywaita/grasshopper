# GrassHopper

very simply backup solution. GitHub as backend!

## Usage

- Execute binary

```shell script
$ grasshopper you-want-to-backup-file
```

- Automatic `git commit` and `git push`
- You can view backup file on GitHub!

## Configuration
### Embedded binary (Recommend)

- Get GitHub Personal token
  - ref: https://github.blog/2013-05-16-personal-api-tokens/
- build binary with embedded value
  - Be handle binaries carefully, do not publish!

```shell script
// change me to oneline
$ go build . -o grasshopper -ldflags '-X main.DefaultGitHubRepository=<your repo> 
  -X main.DefaultGitHubUser=<your user> 
  -X main.DefaultGitHubToken=<your token> '
$ mv grasshopper /your/binary/path/
```

- Execute!

```shell script
$ grasshopper you-want-to-backup-file

# or

$ grasshopper /etc/you-want-to-backup-file    
```

### Environment value

You can set to config by Environment value.

- `GH_GITHUB_REPO`
- `GH_GITHUB_USER`
- `GH_GITHUB_TOKEN`