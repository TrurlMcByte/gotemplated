package main

import (
	"encoding/json"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"text/template"

	"dario.cat/mergo"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func is_map(i interface{}) bool {
	var tst map[string]interface{}
	return reflect.TypeOf(i) == reflect.TypeOf(tst)
}

func map_have(i map[string]interface{}, name string) bool {
	if _, ok := i[name]; ok {
		return true
	}
	return false
}

func env(name string) string {
	return os.Getenv(name)
}

func envdef(name string, def string) string {
	if ret := os.Getenv(name); ret != "" {
		return ret
	}
	return def
}

func LoadJSON_file(filename string, cmap string) map[string]interface{} {
	var jdata map[string]interface{}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil
	}
	if cmap != "" {
		data = []byte("{\"" + cmap + "\":" + string(data) + "}")
	}
	err = json.Unmarshal(data, &jdata)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil
	}
	return jdata
}

func LoadJSON_string(str string, cmap string) map[string]interface{} {
	var jdata map[string]interface{}
	if cmap != "" {
		str = "{\"" + cmap + "\":" + str + "}"
	}
	err := json.Unmarshal([]byte(str), &jdata)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil
	}
	return jdata
}

func LoadJSON_url(url string, cmap string) map[string]interface{} {
	var jdata map[string]interface{}
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching URL:", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Println(url, " returned ", resp.Status)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}
	if cmap != "" {
		body = []byte("{\"" + cmap + "\":" + string(body) + "}")
	}
	err = json.Unmarshal(body, &jdata)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil
	}
	return jdata
}

func LoadTPL_url(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching URL:", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Println(url, " returned ", resp.Status)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}
	return string(body)
}

func LoadTPL_file(filename string) string {
	tdat, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file:", err)
		return ""
	}
	return string(tdat)
}

func usage() {
	log.Fatalln(`gotemplated usage:
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
      envdef {variable} {default}       return environment variable as string if not empty, overwise return default`)
}

func main() {

	var conf map[string]interface{}
	var cmap string = ""

	funcMap := template.FuncMap{
		"is_map":   is_map,
		"map_have": map_have,
		"env":      env,
		"envdef":   envdef,
	}

	t, err := template.New("config").Funcs(funcMap).Parse("")
	if err != nil {
		log.Fatal(err)
		usage()
	}
	if len(os.Args) == 1 {
		usage()
	}

	dirmode := os.FileMode(int(0775))
	filemode := os.FileMode(int(0664))
	uid := os.Geteuid()
	gid := os.Getgid()

	for a := 1; a < len(os.Args); a += 1 {
		if os.Args[a] == "--help" {
			usage()
		} else if os.Args[a] == "--print" {
			a++
			log.Println(os.Args[a])
		} else if os.Args[a] == "--printconf" {
			jsonString, _ := json.MarshalIndent(conf, " ", "  ")
			log.Println(string(jsonString))

		} else if os.Args[a] == "--jmap" {
			a++
			cmap = html.EscapeString(os.Args[a])
		} else if os.Args[a] == "--tfile" {
			a++
			tpl := LoadTPL_file(os.Args[a])
			t, err = template.New("config").Funcs(funcMap).Parse(tpl)
			check(err)
		} else if os.Args[a] == "--turl" {
			a += 1
			tpl := LoadTPL_url(os.Args[a])
			t, err = template.New("config").Funcs(funcMap).Parse(tpl)
			check(err)
		} else if os.Args[a] == "--jurl" {
			a += 1
			conf1 := LoadJSON_url(os.Args[a], cmap)
			err = mergo.Merge(&conf, conf1)
			check(err)
		} else if os.Args[a] == "--jstr" {
			a += 1
			conf1 := LoadJSON_string(os.Args[a], cmap)
			err = mergo.Merge(&conf, conf1)
			check(err)
		} else if os.Args[a] == "--jfile" {
			a += 1
			conf1 := LoadJSON_file(os.Args[a], cmap)
			err = mergo.Merge(&conf, conf1)
			check(err)
		} else if os.Args[a] == "--odp" {
			a += 1
			perm, err := strconv.ParseUint(os.Args[a], 8, 32)
			check(err)
			dirmode = os.FileMode(perm)
		} else if os.Args[a] == "--ofp" {
			a += 1
			perm, err := strconv.ParseUint(os.Args[a], 8, 32)
			check(err)
			filemode = os.FileMode(perm)
		} else if os.Args[a] == "--uid" {
			a += 1
			uid, err = strconv.Atoi(os.Args[a])
			check(err)
		} else if os.Args[a] == "--gid" {
			a += 1
			gid, err = strconv.Atoi(os.Args[a])
			check(err)
		} else if os.Args[a] == "--ofile" {
			a += 1
			filename := os.Args[a]
			dirname, err := filepath.Abs(filepath.Dir(filename))
			check(err)
			if _, err := os.Stat(dirname); os.IsNotExist(err) {
				log.Println("create path " + dirname)
				err = os.MkdirAll(dirname, dirmode)
				check(err)
				err = os.Chown(dirname, uid, gid)
				check(err)
			}
			fo, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filemode)
			check(err)
			err = t.Execute(fo, conf)
			check(err)
			err = fo.Close()
			check(err)
			err = os.Chown(filename, uid, gid)
			check(err)
		}
	}
}
