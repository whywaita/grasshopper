package file

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

/*
.
├── backup
│   ├── web01.example.com
│       └── etc/nginx.conf
│   └── db01.example.com
│       └── etc/my.cnf
*/

func ToTreePath(fp string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", errors.Wrap(err, "failed to get Hostname")
	}

	p := filepath.Join(hostname, fp)

	return p, nil
}
