# GrassHopper

simply backup solution to GitHub.

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
- Edit config.go
  - Add Repository / User / Personal API Token

```shell script
$ vim ./storage/github_config.go
```

- build binary

```shell script
$ go build . -o grasshopper
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