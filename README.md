# gver

gver is a tool which uses [Conventional commits](https://www.conventionalcommits.org/) 
to determine a [SemVer](https://semver.org/).

## Examples

Get the help:
```bash
gver -h
```

Determine the version of code in the current directory:
```bash
gver
```

Determine the version of code in a git repository:
```bash
gver -repo /path/to/my/repo
```

