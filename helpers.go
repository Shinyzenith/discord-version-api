package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func get_response_data(url string) (string, error) {
	/* Makes a GET request to a given url and returns the response body as a string. */
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	response_body, _ := io.ReadAll(res.Body)
	return string(response_body), nil
}

func parse_release_channel(channel string) (string, error) {
	/* Parse release channel. */
	release_channel := strings.ToLower(strings.TrimSpace(channel))
	if release_channel != "canary" && release_channel != "stable" && release_channel != "ptb" {
		return "", errors.New("Invalid release channel passed")
	}

	/* stable.discord.com does not exist so we simply remove it. */
	if release_channel == "stable" {
		return "", nil
	} else {
		return fmt.Sprintf("%s.", release_channel), nil
	}
}

func search_build_info(url string, go_channel chan string) {
	/* Check for Build info, if we get a match then we send it over the channel. */
	re := regexp.MustCompile(`Build Number: [\w.]+, Version Hash: [\w.]+`)
	response_body, _ := get_response_data(url)
	for _, match := range re.FindAllString(response_body, -1) {
		go_channel <- match
	}
}
