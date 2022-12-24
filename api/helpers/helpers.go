package helpers

import (
	"os"
	"strings"
)

func Enforehttp(url string) string{
	if url[:4] != "http"{
		return "http://" + url
	}

	return url
}


func DomainError(url string)bool{


	if url == os.Getenv("DOMAIN"){
		return true
	}

	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(url, "https://", "", 1)
	newUrl = strings.Replace(url, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]

	if newUrl != os.Getenv("DOMAIN"){
		return false
	}

	return true
}