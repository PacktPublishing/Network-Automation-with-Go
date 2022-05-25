package tool

import (
	"tool/cli"
	"encoding/json"
	"encoding/base64"
	"strings"
	"tool/http"
	"network.automation:cvx"
)

devices: ["clab-netgo-cvx"]
auth: base64.Encode(null, "cumulus:cumulus")

save: state: "apply"
save: "auto-prompt": ays:           "ays_yes"
save: "auto-prompt": "ignore_fail": "ignore_fail_yes"

command: apply: {
	for _, d in devices {
		(d): {

			create: http.Post & {
				url: "https://\(d):8765/nvue_v1/revision"
				tls: verify: false
				request: header: "Authorization": "Basic \(auth)"
			}

			revisionID: [ for k, v in json.Unmarshal(create.response.body) {k}]
			escapedID: strings.Replace(revisionID[0], "/", "%2F", -1)

			patch: http.Do & {
				method: "PATCH"
				url:    "https://\(d):8765/nvue_v1/?rev=\(escapedID)"
				tls: verify: false
				request: header: "Authorization": "Basic \(auth)"
				request: header: "Content-Type":  "application/json"
				request: body: json.Marshal(cvx)
			}

			print: cli.Print & {
				text: "PATCH RESPONSE: \(patch.response.body)"
			}

			apply: http.Do & {
				$after: patch
				method: "PATCH"
				url:    "https://\(d):8765/nvue_v1/revision/\(escapedID)"
				tls: verify: false
				request: header: "Authorization": "Basic \(auth)"
				request: header: "Content-Type":  "application/json"
				request: body: json.Marshal(save)
			}

			ok: cli.Print & {
				text: "APPLY RESPONSE \(apply.response.body)"
			}
		}
	}
}
