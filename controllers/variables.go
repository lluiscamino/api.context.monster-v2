package controllers

import (
	"os"
	"strconv"
)

var maxLimit, _ = strconv.Atoi(os.Getenv("max_limit"))
