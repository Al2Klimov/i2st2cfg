package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"i2st2cfg/lib/base"
	"io"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

Objects:
	for {
		ns, errRN := base.ReadNetStringFromStream(in, -1)
		if errRN == io.EOF {
			break
		}

		var obj struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}
		if errUJ := json.Unmarshal(ns, &obj); errUJ != nil {
			fmt.Fprintln(os.Stderr, errUJ.Error())
			os.Exit(1)
		}

		nameParts := strings.Split(obj.Name, "!")

		switch obj.Type {
		case "ApiListener", "Endpoint", "LivestatusListener":
			continue Objects
		case "CheckCommand":
			switch obj.Name {
			case "ido":
				continue Objects
			}
		}

		fmt.Fprintf(out, "object Types[%s] %s {\n", renderString(obj.Type), renderString(nameParts[len(nameParts)-1]))

		switch obj.Type {
		case "Comment":
			fmt.Fprintln(out, `  author = " "`)
			fmt.Fprintln(out, `  host_name = `+renderString(nameParts[0]))

			if len(nameParts) > 2 {
				fmt.Fprintln(out, `  service_name = `+renderString(nameParts[1]))
			}

			fmt.Fprintln(out, `  text = " "`)
		case "Downtime":
			fmt.Fprintln(out, `  author = " "`)
			fmt.Fprintln(out, `  comment = " "`)
			fmt.Fprintln(out, `  host_name = `+renderString(nameParts[0]))

			if len(nameParts) > 2 {
				fmt.Fprintln(out, `  service_name = `+renderString(nameParts[1]))
			}

			fmt.Fprintln(out, `  start_time = get_time()`)
			fmt.Fprintln(out, `  end_time = start_time + 30m`)
		case "FileLogger":
			fmt.Fprintln(out, `  path = "/dev/null"`)
			fmt.Fprintln(out, `  severity = "information"`)
		case "Host":
			fmt.Fprintln(out, `  check_command = "ido"`)
		case "Notification":
			fmt.Fprintln(out, "  var nc = get_objects(NotificationCommand)")
			fmt.Fprintln(out, `  command = nc[random() % len(nc)].name`)
			fmt.Fprintln(out, `  host_name = `+renderString(nameParts[0]))

			if len(nameParts) > 2 {
				fmt.Fprintln(out, `  service_name = `+renderString(nameParts[1]))
			}

			fmt.Fprintln(out, `  user_groups = get_objects(UserGroup).map(ug => ug.name)`)
			//fmt.Fprintln(out, `  users = get_objects(User).map(u => u.name`)
		case "Service":
			fmt.Fprintln(out, "  var cc = get_objects(CheckCommand)")
			fmt.Fprintln(out, `  check_command = cc[random() % len(cc)].name`)
			fmt.Fprintln(out, `  host_name = `+renderString(nameParts[0]))
		}

		fmt.Fprintln(out, "}")
	}

	out.Flush()
}

func renderString(s string) string {
	return "{{{" + strings.Replace(s, "}}}", `}}}+"}}}"+{{{`, -1) + "}}}"
}
