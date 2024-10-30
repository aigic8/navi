# Navi

Go CLI Application to tag files and search through them.

## Basic Idea

The idea is to create tags and save them in a config file.

- The config file will be local, this way if the files are put in a git repository the config file will also be shared.

### Structure of the Config file

```yaml
files:
  papers.review.r001_how-to-become-a-researcher.pdf:
    - review
    - gold
    - needs-review
  papers.modeling.**:
    - modeling
```

- `*` means all the files in the directory, `**` means all the files in the directory and subdirectory
  - JUST AN IDEA FOR SIMPLICITY IN WRITING

### Location and Name of the Config

The config file would be `navi.yaml` or `navi.yml` in the root directory of the project

### Commands

#### Finding a file based on tag

```bash
navi find -t gold
```

an alternative (for faster typing) would be

```bash
navi f -t gold
```

Having an `OR` filter:

```bash
navi f -t gold|modeling
```

Having an `AND` filter:

```bash
navi f -t gold&modeling
```

Combining both:

```bash
navi f -t gold&(modeling|optimization)
```

#### Adding tags to files/dirs

```bash
navi add . -t readLater
```

Or a shorter syntax for `add` would be:

```bash
navi a . -t readLater
```

#### Getting a list of files in a dir

```bash
navi ls ./dir
```

This would result a similar text to below:

```
./modeling/basics.of.modeling.pdf   modeling,gold
./modeling/finding.important.parameters.in.models.pdf modeling,readLater
```
