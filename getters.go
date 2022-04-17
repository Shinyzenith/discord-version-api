package main

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

type BuildInfo struct {
	BuildNumber string
	VersionHash string
}

var collector = colly.NewCollector() // Global colly collector.

func init() {
	collector.AllowURLRevisit = true
}

func get_build_info(channel string) (BuildInfo, error) {
	release_channel, err := parse_release_channel(channel)
	if err != nil {
		return BuildInfo{}, err
	}

	/* Get the page and parse it's body into string from a stream of bytes. */
	response_body, err := get_response_data(fmt.Sprintf("https://%sdiscord.com/app", release_channel))
	if err != nil {
		return BuildInfo{}, err
	}

	/* Make a channel for the goroutines to communicate over. */
	go_channel := make(chan string)

	/* Get all matching js files and concurrently check them all for version info. */
	re := regexp.MustCompile(`\/assets\/\w*\.js`)
	for _, match := range re.FindAllString(response_body, -1) {
		go search_build_info(fmt.Sprintf("https://%sdiscord.com%s", release_channel, match), go_channel)
	}

	/* If we get a match then we simply parse it and return it. */
	for msg := range go_channel {

		/* Extracting the data. */
		build_data := BuildInfo{}
		for index, element := range strings.Split(msg, ",") {
			split := strings.Split(element, " ")
			switch index {
			case 0:
				build_data.BuildNumber = split[len(split)-1]
			case 1:
				build_data.VersionHash = split[len(split)-1]
			}
		}

		return build_data, nil
	}

	/* If we get nothing then we error out. */
	return BuildInfo{}, errors.New("Could not find any version info.")
}

func get_build_id(channel string) (string, error) {
	/* Parse release channel. */
	release_channel, err := parse_release_channel(channel)
	if err != nil {
		return "", err
	}

	/* Make a GET request. */
	env_var := ""
	collector.OnHTML("script[nonce]", func(element *colly.HTMLElement) {
		if strings.Contains(element.Text, "SENTRY_TAGS") {
			env_var = element.Text
		}
	})

	collector.Visit(fmt.Sprintf("https://%sdiscord.com/app", release_channel))

	if env_var != "" {
		for _, line := range strings.Split(env_var, ",") {
			/* Parse the buildId */
			if !strings.Contains(line, "SENTRY_TAGS") || !strings.Contains(line, "buildId") {
				continue
			}
			re := regexp.MustCompile(`^\s*SENTRY_TAGS: \{"buildId":"(.*)"$`)
			return re.ReplaceAllString(line, "${1}"), nil
		}
	}

	return "", errors.New("Couldn't find environment variables.")
}

func get_android_stable_id() (string, error) {
	re := regexp.MustCompile(`(\d*\.\d*) - Stable`)
	match := ""
	collector.OnHTML("span.htlgb", func(element *colly.HTMLElement) {
		res := re.FindAllString(element.Text, -1)
		if len(res) == 0 {
			return
		}
		match = re.ReplaceAllString(element.Text, "${1}")
	})
	collector.Visit("https://play.google.com/store/apps/details?id=com.discord")
	if match == "" {
		return "", errors.New("Couldn't find android stable channel version info")
	} else {
		return match, nil
	}
}
