package aiservices

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractGitHubInfo(url string) (map[string]string, error) {
	if !strings.Contains(url, "github.com") {
		return nil, errors.New("not a valid GitHub URL")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch GitHub profile")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	profile := make(map[string]string)

	doc.Find("div.p-note.user-profile-bio").Each(func(i int, s *goquery.Selection) {
		profile["bio"] = strings.TrimSpace(s.Text())
	})

	doc.Find("a[href$='followers']").Each(func(i int, s *goquery.Selection) {
		profile["followers"] = strings.TrimSpace(s.Text())
	})

	doc.Find("a[href$='?tab=repositories']").Each(func(i int, s *goquery.Selection) {
		profile["repositories"] = strings.TrimSpace(s.Text())
	})

	return profile, nil
}

func ExtractLinkedInInfo(url string) (map[string]string, error) {
	if !strings.Contains(url, "linkedin.com/in/") {
		return nil, errors.New("not a valid LinkedIn profile URL")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch LinkedIn profile")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	profile := make(map[string]string)

	doc.Find("h1.text-heading-xlarge").Each(func(i int, s *goquery.Selection) {
		profile["name"] = strings.TrimSpace(s.Text())
	})

	doc.Find("div.text-body-medium.break-words").Each(func(i int, s *goquery.Selection) {
		profile["headline"] = strings.TrimSpace(s.Text())
	})

	return profile, nil
}
