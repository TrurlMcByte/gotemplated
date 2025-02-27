[![Build Status](https://travis-ci.org/TrurlMcByte/gotemplated.svg?branch=master)](https://travis-ci.org/TrurlMcByte/gotemplated)

# gotemplated

Go based "templated" tools
Generate config files (but not only) from templates (include url) and json data sources (etcd, firebase, webdis, etc...)

# get

ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, stripped available from TravisCI builds on page <https://github.com/TrurlMcByte/gotemplated/releases>

# build

```bash
go install -gccgoflags="-w -s" github.com/TrurlMcByte/gotemplated@latest
```

for static build

```bash
CGO_ENABLED=0  go build -gccgoflags "-s"
```

# documentation

all parametrs may be called multiple times and used in order of presense

```text
    --jurl {URL}        load and merge data from URL (JSON)
    --jfile {FILE}      load and merge data from file (JSON)
    --jstr {STRING}     load and merge data from string (JSON)
    --tfile {FILE}      load template from file
    --turl {URL}        load template from URL
    --odp {PERM}        default permissions for created directories (octal)
    --ofp {PERM}        default permissions for created files (octal)
    --uid {UID}         default owner (uid) for created files (int)
    --gid {GID}         default group (gid) for created files (int)
    --ofile {FILE}      execute last template and write result to file (also create path if not exists)
    --jmap {STRING}     map next loaded json string/file/url to subject, empty string reset mapping
    --print {ARG}       print argument (for debug)
    --printconf         print all collected data (for debug)

    Additional template funtions:
      is_map {variable}                 return true if variable is map
      map_have {variable} "string"      return true if variable have field "string"
      env {variable}                    return environment variable as string
      envdef {variable} {default}       return environment variable as string if not empty, overwise return default

```

see also <https://golang.org/pkg/text/template/> for templates

I recommend to use "jq" or "yq" tools to prepare data

# usage

```bash
./gotemplated \
    --print "Load data from url" \
    --jurl https://myconf.firebaseio.com/conf/test.json \
    --print "Load data from file and merge with prev" \
    --jfile meta.json \
    --print "Load data from string and merge with prev" \
    --jstr '{"ctrl": "Some additional data"}' \
    --print "Load template #1 from url" \
    --turl https://myconf.firebaseapp.com/test.tpl \
    --print "Execute template #1 and write to file" \
    --ofile test1.res \
    --print "Load template #2 from file" \
    --tfile test.tpl \
    --odp 0700 \
    --ofp 0400 \
    --uid 82 \
    --gid 82 \
    --print "Create /tmp/data with permissions 0700/82.82 and execute template #2 and write to file with permissions 0400/82.82" \
    --ofile /tmp/data/test2.res
```

# example

```bash
gotemplated.exe --jstr "{\"hostname\":\"myhost.example.org\"}" --tfile some.tpl --ofile some.conf
```

where ```some.tpl``` containts

```text
hostname = "{{ .hostname }}";
home = "{{ envdef "HOME" "/none" }}"
```

in result in file ```some.conf``` will be

```ini
hostname = "myhost.example.org"
```
