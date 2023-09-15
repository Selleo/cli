package random

import "github.com/gofrs/uuid/v5"

func UUID4() string {
	u := uuid.Must(uuid.NewV4())
	return u.String()
}
