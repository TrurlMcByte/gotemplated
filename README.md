# gotemplated
Go based "templated" tools
Generate config files (but not only) from templates (include url) and json data sources (etcd, firebase, webdis, etc...)


# build
go build

# manual
all parametrs may be called multiple times and user in order of presense
```
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
```

# usage
```
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

